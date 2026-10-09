package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/service"
)

type AuthHandler struct {
	authService  *service.AuthService
	authMW       *middleware.AuthMiddleware
	adminHandler *AdminHandler
}

func NewAuthHandler(authService *service.AuthService, authMW *middleware.AuthMiddleware, adminHandlers ...*AdminHandler) *AuthHandler {
	h := &AuthHandler{
		authService: authService,
		authMW:      authMW,
	}
	if len(adminHandlers) > 0 {
		h.adminHandler = adminHandlers[0]
	}
	return h
}

type RegisterRequest struct {
	DisplayName     string `json:"displayName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	AcceptedTerms   bool   `json:"acceptedTerms"`
}

type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type PasswordResetRequest struct {
	Email string `json:"email"`
}

type ConfirmPasswordResetRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ.")
		return
	}

	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.Email = strings.TrimSpace(req.Email)

	if req.DisplayName == "" {
		RespondError(w, http.StatusBadRequest, "Tên hiển thị không được để trống.")
		return
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		RespondError(w, http.StatusBadRequest, "Địa chỉ email không hợp lệ.")
		return
	}
	if len(req.Password) < 8 {
		RespondError(w, http.StatusBadRequest, "Mật khẩu phải chứa ít nhất 8 ký tự.")
		return
	}
	if req.Password != req.ConfirmPassword {
		RespondError(w, http.StatusBadRequest, "Mật khẩu xác nhận không khớp.")
		return
	}
	if !req.AcceptedTerms {
		RespondError(w, http.StatusBadRequest, "Bạn cần đồng ý với Điều khoản sử dụng.")
		return
	}

	accessToken, refreshToken, user, err := h.authService.Register(
		r.Context(),
		req.DisplayName,
		req.Email,
		req.Password,
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			RespondError(w, http.StatusConflict, "Email này đã được sử dụng. Vui lòng đăng nhập hoặc sử dụng email khác.")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Không thể đăng ký tài khoản. Vui lòng thử lại sau.")
		return
	}

	h.authMW.SetAuthCookies(w, accessToken, refreshToken)

	RespondJSON(w, http.StatusOK, domain.AuthSessionResponse{
		User: domain.UserDTO{
			ID:          user.ID,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			Role:        user.Role,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.login(w, r, false)
}

func (h *AuthHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	h.login(w, r, true)
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request, staffOnly bool) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ.")
		return
	}

	req.Identity = strings.TrimSpace(req.Identity)
	if req.Identity == "" || req.Password == "" {
		RespondError(w, http.StatusBadRequest, "Vui lòng nhập đầy đủ thông tin đăng nhập.")
		return
	}

	login := h.authService.Login
	if staffOnly {
		login = h.authService.LoginStaff
	}
	accessToken, refreshToken, user, err := login(
		r.Context(),
		req.Identity,
		req.Password,
		req.Remember,
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		if errors.Is(err, service.ErrStaffRequired) {
			RespondError(w, http.StatusForbidden, "Tài khoản này không có quyền quản trị.")
			return
		}
		RespondError(w, http.StatusUnauthorized, "Email hoặc mật khẩu không chính xác.")
		return
	}

	h.authMW.SetAuthCookies(w, accessToken, refreshToken)

	RespondJSON(w, http.StatusOK, domain.AuthSessionResponse{
		User: domain.UserDTO{
			ID:          user.ID,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			Role:        user.Role,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		_ = h.authService.Logout(r.Context(), cookie.Value)
	}

	h.authMW.ClearAuthCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa đăng nhập.")
		return
	}

	RespondJSON(w, http.StatusOK, domain.AuthSessionResponse{
		User: domain.UserDTO{
			ID:          userCtx.ID,
			DisplayName: userCtx.DisplayName,
			Email:       userCtx.Email,
			Role:        userCtx.Role,
		},
	})
}

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ.")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		RespondError(w, http.StatusBadRequest, "Vui lòng nhập một địa chỉ email hợp lệ.")
		return
	}

	_ = h.authService.RequestPasswordReset(r.Context(), req.Email)

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) ValidatePasswordReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ.")
		return
	}
	if err := h.authService.ValidatePasswordReset(r.Context(), req.Token); err != nil {
		respondPasswordResetError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondPasswordResetError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrPasswordResetNotFound) {
		RespondError(w, http.StatusBadRequest, "Liên kết đặt lại mật khẩu không hợp lệ, đã được sử dụng hoặc đã hết hạn. Vui lòng yêu cầu liên kết mới.")
		return
	}
	RespondError(w, http.StatusInternalServerError, "Không thể xử lý yêu cầu đặt lại mật khẩu. Vui lòng thử lại sau.")
}

func (h *AuthHandler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	var req ConfirmPasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ.")
		return
	}

	if req.Token == "" || len(req.NewPassword) < 8 {
		RespondError(w, http.StatusBadRequest, "Thông tin đặt lại mật khẩu không hợp lệ.")
		return
	}

	if err := h.authService.ConfirmPasswordReset(r.Context(), req.Token, req.NewPassword); err != nil {
		respondPasswordResetError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Mật khẩu đã được cập nhật thành công. Vui lòng đăng nhập bằng mật khẩu mới.",
	})
}
