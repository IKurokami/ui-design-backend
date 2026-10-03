package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"ui-design-backend/internal/auth"
	"ui-design-backend/internal/config"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/mailer"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
	"ui-design-backend/internal/service"
	transporthttp "ui-design-backend/internal/transport/http"
)

func setupTestApp(t *testing.T) (http.Handler, *auth.TokenService, func()) {
	t.Helper()

	dbFile := "./test_auth.db"
	_ = os.Remove(dbFile)

	cfg := &config.Config{
		Port:                    "8080",
		Env:                     "development",
		ClientOrigin:            "http://localhost:5173",
		DBDriver:                "sqlite",
		DBDSN:                   dbFile,
		JWTSecret:               "test-secret-key-1234567890",
		AccessTokenTTL:          1 * time.Second, // short TTL for test
		RefreshTokenTTL:         10 * time.Minute,
		RefreshTokenRememberTTL: 30 * time.Minute,
		RotationGracePeriod:     5 * time.Second,
		SMTPHost:                "127.0.0.1",
		SMTPPort:                1025,
		MailFrom:                "test@uidesign.local",
		AppBaseURL:              "http://localhost:5173",
	}

	db, err := repository.NewDB(cfg.DBDriver, cfg.DBDSN)
	if err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}

	userRepo := repository.NewUserRepo(db)
	sessionRepo := repository.NewSessionRepo(db)
	resetRepo := repository.NewPasswordResetRepo(db)

	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL)
	mailService := mailer.NewMailer(cfg)
	authService := service.NewAuthService(cfg, userRepo, sessionRepo, resetRepo, tokenService, mailService)

	authMW := middleware.NewAuthMiddleware(cfg, tokenService, authService)
	authHandler := transporthttp.NewAuthHandler(authService, authMW)

	router := transporthttp.NewRouter(cfg, authHandler, authMW)

	cleanup := func() {
		_ = db.Close()
		_ = os.Remove(dbFile)
	}

	return router, tokenService, cleanup
}

func TestTransparentAutoRefresh(t *testing.T) {
	router, tokenService, cleanup := setupTestApp(t)
	defer cleanup()

	// 1. Register user
	registerBody := `{"displayName":"Tran Van B","email":"tranvanb@example.com","password":"Password123","confirmPassword":"Password123","acceptedTerms":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(registerBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Register failed: status %d, body: %s", w.Code, w.Body.String())
	}

	resCookies := w.Result().Cookies()
	var accessTokenCookie, refreshTokenCookie *http.Cookie
	for _, c := range resCookies {
		if c.Name == "access_token" {
			accessTokenCookie = c
		} else if c.Name == "refresh_token" {
			refreshTokenCookie = c
		}
	}

	if accessTokenCookie == nil || refreshTokenCookie == nil {
		t.Fatalf("Expected both access_token and refresh_token cookies to be set on register")
	}

	// 2. Call GET /api/auth/me immediately while access_token is VALID
	reqMe := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	reqMe.AddCookie(accessTokenCookie)
	wMe := httptest.NewRecorder()
	router.ServeHTTP(wMe, reqMe)

	if wMe.Code != http.StatusOK {
		t.Fatalf("GetMe with valid access_token failed: %d, body: %s", wMe.Code, wMe.Body.String())
	}

	// 3. Now intentionally WAIT for access token to expire (AccessTokenTTL was set to 1 second)
	time.Sleep(1200 * time.Millisecond)

	// Verify the access token is indeed expired now
	_, err := tokenService.ValidateAccessToken(accessTokenCookie.Value)
	if err == nil {
		t.Fatalf("Expected access token to be expired after sleep")
	}

	// 4. NOW THE CORE TEST: Client sends GET /api/auth/me with:
	// - Expired access_token cookie
	// - Valid refresh_token cookie
	// Client DOES NOT send a separate refresh request!
	reqMeRefreshed := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	reqMeRefreshed.AddCookie(accessTokenCookie)  // EXPIRED
	reqMeRefreshed.AddCookie(refreshTokenCookie) // VALID

	wMeRefreshed := httptest.NewRecorder()
	router.ServeHTTP(wMeRefreshed, reqMeRefreshed)

	// MUST be HTTP 200! ZERO RETRY NEEDED!
	if wMeRefreshed.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200 with transparent auto-refresh, got %d: %s", wMeRefreshed.Code, wMeRefreshed.Body.String())
	}

	var sessionResp domain.AuthSessionResponse
	if err := json.Unmarshal(wMeRefreshed.Body.Bytes(), &sessionResp); err != nil {
		t.Fatalf("Failed to decode session response: %v", err)
	}

	if sessionResp.User.Email != "tranvanb@example.com" {
		t.Errorf("Expected user email tranvanb@example.com, got %s", sessionResp.User.Email)
	}

	// VERIFY Set-Cookie headers in the response containing NEW access_token and rotated refresh_token!
	refreshedCookies := wMeRefreshed.Result().Cookies()
	var newAccessTokenCookie, newRefreshTokenCookie *http.Cookie
	for _, c := range refreshedCookies {
		if c.Name == "access_token" {
			newAccessTokenCookie = c
		} else if c.Name == "refresh_token" {
			newRefreshTokenCookie = c
		}
	}

	if newAccessTokenCookie == nil {
		t.Fatalf("Expected new access_token cookie to be returned in Set-Cookie header")
	}
	if newAccessTokenCookie.Value == accessTokenCookie.Value {
		t.Errorf("Expected new access_token value to differ from old expired one")
	}
	if newRefreshTokenCookie == nil {
		t.Fatalf("Expected rotated refresh_token cookie in Set-Cookie header")
	}

	// 5. Test Grace Period: A concurrent in-flight request arrives with the OLD refresh token within 5s
	reqConcurrent := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	reqConcurrent.AddCookie(accessTokenCookie)  // old expired
	reqConcurrent.AddCookie(refreshTokenCookie) // old refresh token, now inside grace period

	wConcurrent := httptest.NewRecorder()
	router.ServeHTTP(wConcurrent, reqConcurrent)

	if wConcurrent.Code != http.StatusOK {
		t.Fatalf("Expected concurrent in-flight request within grace period to succeed with 200, got %d: %s", wConcurrent.Code, wConcurrent.Body.String())
	}

	t.Logf("SUCCESS: Transparent auto-refresh seamlessly returned HTTP 200 and updated cookies!")
}

func TestPasswordResetFlow(t *testing.T) {
	router, _, cleanup := setupTestApp(t)
	defer cleanup()

	// 1. Register user
	registerBody := `{"displayName":"Nguyen Van C","email":"nguyenvanc@example.com","password":"OldPassword123","confirmPassword":"OldPassword123","acceptedTerms":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(registerBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Register failed: %d", w.Code)
	}

	// 2. Request password reset
	resetReqBody := `{"email":"nguyenvanc@example.com"}`
	reqReset := httptest.NewRequest(http.MethodPost, "/api/auth/password-reset", bytes.NewBufferString(resetReqBody))
	reqReset.Header.Set("Content-Type", "application/json")
	wReset := httptest.NewRecorder()
	router.ServeHTTP(wReset, reqReset)

	if wReset.Code != http.StatusNoContent {
		t.Fatalf("Expected 204 No Content for password reset request, got %d", wReset.Code)
	}

	// Allow goroutine to send mail to Mailpit
	time.Sleep(300 * time.Millisecond)

	// 3. Request reset with non-existent email (must also return 204 to avoid enumeration)
	nonExistentBody := `{"email":"random-not-found@example.com"}`
	reqNonExistent := httptest.NewRequest(http.MethodPost, "/api/auth/password-reset", bytes.NewBufferString(nonExistentBody))
	reqNonExistent.Header.Set("Content-Type", "application/json")
	wNonExistent := httptest.NewRecorder()
	router.ServeHTTP(wNonExistent, reqNonExistent)

	if wNonExistent.Code != http.StatusNoContent {
		t.Fatalf("Expected 204 No Content for non-existent email, got %d", wNonExistent.Code)
	}
}

func TestUnauthenticatedRequest(t *testing.T) {
	router, _, cleanup := setupTestApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Phiên đăng nhập đã hết hạn") {
		t.Errorf("Unexpected error body: %s", w.Body.String())
	}
}
