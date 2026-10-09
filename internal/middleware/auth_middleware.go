package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"ui-design-backend/internal/auth"
	"ui-design-backend/internal/config"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/service"
)

type contextKey string

const UserContextKey contextKey = "auth_user"

type AuthContextUser struct {
	ID          string      `json:"id"`
	Email       string      `json:"email"`
	DisplayName string      `json:"displayName"`
	Role        domain.Role `json:"role,omitempty"`
}

type AuthMiddleware struct {
	cfg          *config.Config
	tokenService *auth.TokenService
	authService  *service.AuthService
	userRepo     domain.UserRepository
}

func NewAuthMiddleware(cfg *config.Config, tokenService *auth.TokenService, authService *service.AuthService, userRepos ...domain.UserRepository) *AuthMiddleware {
	m := &AuthMiddleware{
		cfg:          cfg,
		tokenService: tokenService,
		authService:  authService,
	}
	if len(userRepos) > 0 {
		m.userRepo = userRepos[0]
	}
	return m
}

func (m *AuthMiddleware) currentUser(ctx context.Context, id string) (AuthContextUser, error) {
	var user *domain.User
	var err error
	if m.userRepo != nil {
		user, err = m.userRepo.FindByID(ctx, id)
	} else if m.authService != nil {
		user, err = m.authService.GetUserByID(ctx, id)
	} else {
		return AuthContextUser{}, service.ErrSessionInvalid
	}
	if err != nil {
		return AuthContextUser{}, err
	}
	return AuthContextUser{ID: user.ID, Email: user.Email, DisplayName: user.DisplayName, Role: user.Role}, nil
}

func RequireStaff(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		user, ok := GetUserFromContext(r.Context())
		if !ok {
			respondJSONError(w, http.StatusUnauthorized, "Chưa đăng nhập.")
			return
		}
		if !user.Role.IsStaff() {
			respondJSONError(w, http.StatusForbidden, "Tài khoản này không có quyền quản trị.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r.Context())
		if !ok {
			respondJSONError(w, http.StatusUnauthorized, "Chưa đăng nhập.")
			return
		}
		if user.Role != domain.RoleAdmin {
			respondJSONError(w, http.StatusForbidden, "Chỉ admin được quản trị người dùng.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAuth implements transparent auto-refresh:
func (m *AuthMiddleware) extractAccessToken(r *http.Request) string {
	if authHeader := r.Header.Get("Authorization"); len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	if cookie, err := r.Cookie("access_token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return ""
}

// RequireAuth implements transparent auto-refresh:
// If access token is valid, it proceeds.
// If access token is expired/missing but refresh token is valid,
// it rotates tokens in-flight, attaches Set-Cookie headers, and proceeds to the handler!
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Inspect access token (header or cookie)
		tokenStr := m.extractAccessToken(r)
		if tokenStr != "" {
			claims, err := m.tokenService.ValidateAccessToken(tokenStr)
			if err == nil {
				// Reload permissions and account existence rather than trusting stale JWT data.
				userCtx, userErr := m.currentUser(r.Context(), claims.UserID)
				if userErr != nil {
					m.ClearAuthCookies(w)
					respondJSONError(w, http.StatusUnauthorized, "Phiên đăng nhập không còn hợp lệ. Vui lòng đăng nhập lại.")
					return
				}
				ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// If token was rejected for reasons OTHER than expiration, return 401
			if !errors.Is(err, auth.ErrTokenExpired) {
				m.ClearAuthCookies(w)
				respondJSONError(w, http.StatusUnauthorized, "Mã xác thực không hợp lệ. Vui lòng đăng nhập lại.")
				return
			}
		}

		// 2. Access token is missing or expired: inspect refresh_token cookie
		refreshCookie, err := r.Cookie("refresh_token")
		if err != nil || refreshCookie.Value == "" || m.authService == nil {
			respondJSONError(w, http.StatusUnauthorized, "Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.")
			return
		}

		// 3. Perform transparent in-flight token refresh
		newAccessToken, newRefreshToken, user, err := m.authService.RotateSession(
			r.Context(),
			refreshCookie.Value,
			r.UserAgent(),
			r.RemoteAddr,
		)
		if err != nil {
			m.ClearAuthCookies(w)
			respondJSONError(w, http.StatusUnauthorized, "Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.")
			return
		}

		// 4. Attach fresh cookies to the HTTP response
		m.SetAuthCookies(w, newAccessToken, newRefreshToken)

		// 5. Inject refreshed user into context and execute the actual handler seamlessly!
		userCtx := AuthContextUser{
			ID:          user.ID,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			Role:        user.Role,
		}
		ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth attempts to authenticate the user if cookies or headers are present,
// but allows unauthenticated requests to proceed without error.
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := m.extractAccessToken(r)
		if tokenStr != "" {
			claims, err := m.tokenService.ValidateAccessToken(tokenStr)
			if err == nil {
				userCtx, userErr := m.currentUser(r.Context(), claims.UserID)
				if userErr != nil {
					next.ServeHTTP(w, r)
					return
				}
				ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		refreshCookie, err := r.Cookie("refresh_token")
		if err == nil && refreshCookie.Value != "" && m.authService != nil {
			newAccessToken, newRefreshToken, user, err := m.authService.RotateSession(
				r.Context(),
				refreshCookie.Value,
				r.UserAgent(),
				r.RemoteAddr,
			)
			if err == nil {
				m.SetAuthCookies(w, newAccessToken, newRefreshToken)
				userCtx := AuthContextUser{
					ID:          user.ID,
					Email:       user.Email,
					DisplayName: user.DisplayName,
					Role:        user.Role,
				}
				ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	isSecure := m.cfg.IsProduction()

	// Access Token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(m.cfg.AccessTokenTTL.Seconds()),
	})

	// Refresh Token cookie (only if a new one was provided)
	if refreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   isSecure,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(m.cfg.RefreshTokenRememberTTL.Seconds()),
		})
	}
}

func (m *AuthMiddleware) ClearAuthCookies(w http.ResponseWriter) {
	isSecure := m.cfg.IsProduction()

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func GetUserFromContext(ctx context.Context) (AuthContextUser, bool) {
	user, ok := ctx.Value(UserContextKey).(AuthContextUser)
	return user, ok
}

func respondJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":   http.StatusText(status),
		"message": message,
	})
}
