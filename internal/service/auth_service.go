package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"ui-design-backend/internal/auth"
	"ui-design-backend/internal/config"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/mailer"
	"ui-design-backend/internal/repository"
)

var (
	ErrUserAlreadyExists     = errors.New("email đã được sử dụng")
	ErrInvalidCredentials    = errors.New("thông tin đăng nhập không chính xác")
	ErrStaffRequired         = errors.New("tài khoản không có quyền quản trị")
	ErrSessionInvalid        = errors.New("phiên đăng nhập không hợp lệ")
	ErrSessionExpired        = errors.New("phiên đăng nhập đã hết hạn")
	ErrSessionReuseDetected  = errors.New("phát hiện token bị tái sử dụng bất thường")
	ErrPasswordResetNotFound = errors.New("liên kết đặt lại mật khẩu không hợp lệ, đã được sử dụng hoặc đã hết hạn")
)

type AuthService struct {
	cfg          *config.Config
	userRepo     domain.UserRepository
	sessionRepo  domain.SessionRepository
	resetRepo    domain.PasswordResetRepository
	tokenService *auth.TokenService
	mailer       *mailer.Mailer
	mu           sync.Mutex
}

func NewAuthService(
	cfg *config.Config,
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	resetRepo domain.PasswordResetRepository,
	tokenService *auth.TokenService,
	m *mailer.Mailer,
) *AuthService {
	return &AuthService{
		cfg:          cfg,
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		resetRepo:    resetRepo,
		tokenService: tokenService,
		mailer:       m,
	}
}

func (s *AuthService) Register(ctx context.Context, displayName, email, password string, userAgent, ip string) (string, string, *domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	displayName = strings.TrimSpace(displayName)

	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return "", "", nil, ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return "", "", nil, err
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "", "", nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New().String(),
		DisplayName:  displayName,
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return "", "", nil, err
	}

	accessToken, refreshToken, err := s.createSession(ctx, user, s.cfg.RefreshTokenTTL, userAgent, ip)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) Login(ctx context.Context, identity, password string, remember bool, userAgent, ip string) (string, string, *domain.User, error) {
	return s.login(ctx, identity, password, remember, userAgent, ip, false)
}

func (s *AuthService) LoginStaff(ctx context.Context, identity, password string, remember bool, userAgent, ip string) (string, string, *domain.User, error) {
	return s.login(ctx, identity, password, remember, userAgent, ip, true)
}

func (s *AuthService) login(ctx context.Context, identity, password string, remember bool, userAgent, ip string, staffOnly bool) (string, string, *domain.User, error) {
	identity = strings.ToLower(strings.TrimSpace(identity))

	user, err := s.userRepo.FindByEmail(ctx, identity)
	if err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		return "", "", nil, ErrInvalidCredentials
	}
	if staffOnly && !user.Role.IsStaff() {
		return "", "", nil, ErrStaffRequired
	}

	refreshTTL := s.cfg.RefreshTokenTTL
	if remember {
		refreshTTL = s.cfg.RefreshTokenRememberTTL
	}

	accessToken, refreshToken, err := s.createSession(ctx, user, refreshTTL, userAgent, ip)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) RotateSession(ctx context.Context, rawRefreshToken, userAgent, ip string) (string, string, *domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tokenHash := auth.HashToken(rawRefreshToken)
	session, err := s.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return "", "", nil, ErrSessionInvalid
		}
		return "", "", nil, err
	}

	// 1. Check if session was revoked explicitly
	if session.IsRevoked() {
		return "", "", nil, ErrSessionInvalid
	}

	// 2. Check if already rotated
	if session.RotatedToID != nil {
		// Concurrent request inside grace period!
		if session.InGracePeriod() {
			rotatedSession, err := s.sessionRepo.FindByID(ctx, *session.RotatedToID)
			if err != nil {
				return "", "", nil, ErrSessionInvalid
			}
			user, err := s.userRepo.FindByID(ctx, rotatedSession.UserID)
			if err != nil {
				return "", "", nil, err
			}

			// Issue a fresh access token for the in-flight concurrent request
			accessToken, _, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, user.DisplayName, rotatedSession.ID)
			if err != nil {
				return "", "", nil, err
			}
			// Keep the same rotated session's context, empty string for refresh token so cookie isn't re-mutated needlessly
			return accessToken, "", user, nil
		}

		// Replay attack / Reuse detected outside grace period!
		// Invalidate all sessions for this user for security
		_ = s.sessionRepo.RevokeAllForUser(ctx, session.UserID)
		return "", "", nil, ErrSessionReuseDetected
	}

	// 3. Check if session expired
	if session.IsExpired() {
		return "", "", nil, ErrSessionExpired
	}

	// 4. Regular rotation: create new refresh token and session
	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return "", "", nil, err
	}

	newRawRefreshToken, err := auth.GenerateRandomToken(32)
	if err != nil {
		return "", "", nil, err
	}
	newTokenHash := auth.HashToken(newRawRefreshToken)

	now := time.Now()
	remainingTTL := session.ExpiresAt.Sub(now)
	if remainingTTL < s.cfg.RefreshTokenTTL {
		remainingTTL = s.cfg.RefreshTokenTTL
	}

	newSession := &domain.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: newTokenHash,
		UserAgent: userAgent,
		IP:        ip,
		ExpiresAt: now.Add(remainingTTL),
		CreatedAt: now,
	}

	if err := s.sessionRepo.Rotate(ctx, session.ID, newSession, s.cfg.RotationGracePeriod); err != nil {
		return "", "", nil, err
	}

	newAccessToken, _, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, user.DisplayName, newSession.ID)
	if err != nil {
		return "", "", nil, err
	}

	return newAccessToken, newRawRefreshToken, user, nil
}

func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return nil
	}
	tokenHash := auth.HashToken(rawRefreshToken)
	session, err := s.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil
	}
	return s.sessionRepo.Revoke(ctx, session.ID)
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Printf("[PASSWORD RESET] Email '%s' chưa được đăng ký trong hệ thống. Không gửi email (tiêu chuẩn bảo mật chống rò rỉ tài khoản).", email)
		// Do not reveal whether user exists (prevent email enumeration)
		return nil
	}

	_ = s.resetRepo.InvalidateActiveForUser(ctx, user.ID)

	rawToken, err := auth.GenerateRandomToken(32)
	if err != nil {
		return err
	}
	tokenHash := auth.HashToken(rawToken)

	pr := &domain.PasswordReset{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.resetRepo.Create(ctx, pr); err != nil {
		return err
	}

	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", s.cfg.AppBaseURL, rawToken)

	// Send email asynchronously in goroutine
	go func() {
		if err := s.mailer.SendPasswordReset(user.Email, user.DisplayName, resetURL); err != nil {
			log.Printf("[MAIL ERROR] Failed to send password reset to %s: %v", user.Email, err)
		} else {
			log.Printf("[MAIL SENT] Password reset email sent to %s via Mailpit", user.Email)
		}
	}()

	return nil
}

func (s *AuthService) findValidPasswordReset(ctx context.Context, rawToken string) (*domain.PasswordReset, error) {
	if strings.TrimSpace(rawToken) == "" {
		return nil, ErrPasswordResetNotFound
	}
	pr, err := s.resetRepo.FindByTokenHash(ctx, auth.HashToken(rawToken))
	if errors.Is(err, repository.ErrResetNotFound) {
		return nil, ErrPasswordResetNotFound
	}
	if err != nil {
		return nil, err
	}
	if pr == nil || !pr.IsValid() {
		return nil, ErrPasswordResetNotFound
	}
	return pr, nil
}

// Validation does not consume the token; opening an email link is read-only.
func (s *AuthService) ValidatePasswordReset(ctx context.Context, rawToken string) error {
	_, err := s.findValidPasswordReset(ctx, rawToken)
	return err
}

func (s *AuthService) ConfirmPasswordReset(ctx context.Context, rawToken, newPassword string) error {
	pr, err := s.findValidPasswordReset(ctx, rawToken)
	if err != nil {
		return err
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.resetRepo.Complete(ctx, pr.ID, hashedPassword); err != nil {
		if errors.Is(err, repository.ErrResetNotFound) {
			return ErrPasswordResetNotFound
		}
		return err
	}

	return nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *AuthService) createSession(ctx context.Context, user *domain.User, refreshTTL time.Duration, userAgent, ip string) (string, string, error) {
	rawRefreshToken, err := auth.GenerateRandomToken(32)
	if err != nil {
		return "", "", err
	}
	tokenHash := auth.HashToken(rawRefreshToken)

	now := time.Now()
	session := &domain.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		UserAgent: userAgent,
		IP:        ip,
		ExpiresAt: now.Add(refreshTTL),
		CreatedAt: now,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", "", err
	}

	accessToken, _, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, user.DisplayName, session.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, rawRefreshToken, nil
}
