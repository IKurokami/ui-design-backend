package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
	"ui-design-backend/internal/service"
)

type AdminHandler struct{ service *service.ManagementService }

func NewAdminHandler(s *service.ManagementService) *AdminHandler { return &AdminHandler{service: s} }

type AdminUserDTO struct {
	ID          string      `json:"id"`
	DisplayName string      `json:"displayName"`
	Email       string      `json:"email"`
	Role        domain.Role `json:"role,omitempty"`
	CreatedAt   string      `json:"createdAt"`
}

func userDTO(u *domain.User) AdminUserDTO {
	return AdminUserDTO{ID: u.ID, DisplayName: u.DisplayName, Email: u.Email, Role: u.Role, CreatedAt: u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("paginated") == "true" {
		page, size := parsePage(r)
		users, total, err := h.service.ListUsersPage(r.Context(), page, size, strings.TrimSpace(r.URL.Query().Get("q")), r.URL.Query().Get("role"))
		if err != nil {
			RespondError(w, 500, "Không thể tải người dùng")
			return
		}
		items := make([]AdminUserDTO, 0, len(users))
		for _, u := range users {
			items = append(items, userDTO(u))
		}
		RespondJSON(w, http.StatusOK, pagePayload(items, page, size, total))
		return
	}
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		RespondError(w, 500, "Không thể tải người dùng")
		return
	}
	result := make([]AdminUserDTO, 0, len(users))
	for _, u := range users {
		result = append(result, userDTO(u))
	}
	RespondJSON(w, http.StatusOK, result)
}

func parsePage(r *http.Request) (int, int) {
	p, _ := strconv.Atoi(r.URL.Query().Get("page"))
	n, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if p < 1 {
		p = 1
	}
	if n < 1 {
		n = 20
	}
	if n > 100 {
		n = 100
	}
	return p, n
}
func pagePayload[T any](items []T, page, size, total int) domain.Page[T] {
	pages := total / size
	if total%size != 0 {
		pages++
	}
	return domain.Page[T]{Items: items, Page: page, PageSize: size, Total: total, TotalPages: pages}
}
func (h *AdminHandler) ListStories(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("paginated") == "true" {
		page, size := parsePage(r)
		stories, total, err := h.service.ListStoriesPage(r.Context(), page, size, strings.TrimSpace(r.URL.Query().Get("q")))
		if err != nil {
			RespondError(w, 500, "Không thể tải truyện")
			return
		}
		RespondJSON(w, http.StatusOK, pagePayload(stories, page, size, total))
		return
	}
	stories, err := h.service.ListStories(r.Context())
	if err != nil {
		RespondError(w, 500, "Không thể tải truyện")
		return
	}
	RespondJSON(w, http.StatusOK, stories)
}
func (h *AdminHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Role domain.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, 400, "Dữ liệu quyền không hợp lệ")
		return
	}
	if err := h.service.SetRole(r.Context(), chi.URLParam(r, "id"), body.Role); err != nil {
		if errors.Is(err, service.ErrInvalidRole) {
			RespondError(w, 400, "Role chỉ có thể là admin hoặc moderator")
			return
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			RespondError(w, 404, "Không tìm thấy người dùng")
			return
		}
		RespondError(w, 500, "Không thể cập nhật quyền")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	actor, _ := middleware.GetUserFromContext(r.Context())
	if err := h.service.DeleteUser(r.Context(), actor.ID, chi.URLParam(r, "id")); err != nil {
		if errors.Is(err, service.ErrSelfDelete) {
			RespondError(w, 400, "Không thể xóa tài khoản đang đăng nhập")
			return
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			RespondError(w, 404, "Không tìm thấy người dùng")
			return
		}
		RespondError(w, 500, "Không thể xóa người dùng")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
