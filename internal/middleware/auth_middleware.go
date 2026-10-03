package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"ui-design-backend/internal/auth"
	"ui-design-backend/internal/config"
	"ui-design-backend/internal/service"
)

type contextKey string

const UserContextKey contextKey = "auth_user"

type AuthContextUser struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type AuthMiddleware struct {
	cfg          *config.Config
	tokenService *auth.TokenService
	authService  *service.AuthService
}

func NewAuthMiddleware(cfg *config.Config, tokenService *auth.TokenService, authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:          cfg,
		tokenService: tokenService,
		authService:  authService,
	}
}

// RequireAuth implements transparent auto-refresh:
// If access token is valid, it proceeds.
// If access token is expired/missing but refresh token is valid,
// it rotates tokens in-flight, attaches Set-Cookie headers, and proceeds to the handler!
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Inspect access_token cookie
		accessCookie, err := r.Cookie("access_token")
		if err == nil && accessCookie.Value != "" {
			claims, err := m.tokenService.ValidateAccessToken(accessCookie.Value)
			if err == nil {
				// Access token is valid -> inject user and proceed
				userCtx := AuthContextUser{
					ID:          claims.UserID,
					Email:       claims.Email,
					DisplayName: claims.DisplayName,
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
		if err != nil || refreshCookie.Value == "" {
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
		}
		ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
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
