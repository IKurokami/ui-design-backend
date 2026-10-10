package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
)

func socialUser(w http.ResponseWriter, r *http.Request) (string, bool) {
	w.Header().Set("Cache-Control", "private, no-store")
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa đăng nhập")
		return "", false
	}
	return user.ID, true
}

func (h *StoryHandler) FollowIDs(w http.ResponseWriter, r *http.Request) {
	user, ok := socialUser(w, r)
	if !ok {
		return
	}
	ids, err := h.storyService.FollowIDs(r.Context(), user)
	if err != nil {
		RespondError(w, 500, "Không tải được danh sách theo dõi")
		return
	}
	RespondJSON(w, 200, map[string]any{"ids": ids})
}
func (h *StoryHandler) FollowStory(w http.ResponseWriter, r *http.Request) {
	user, ok := socialUser(w, r)
	if !ok {
		return
	}
	saved, err := h.storyService.Follow(r.Context(), user, chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, 500, "Không thể theo dõi truyện")
		return
	}
	if !saved {
		RespondError(w, 404, "Truyện không thể theo dõi")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *StoryHandler) UnfollowStory(w http.ResponseWriter, r *http.Request) {
	user, ok := socialUser(w, r)
	if !ok {
		return
	}
	if err := h.storyService.Unfollow(r.Context(), user, chi.URLParam(r, "id")); err != nil {
		RespondError(w, 500, "Không thể bỏ theo dõi")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *StoryHandler) Notifications(w http.ResponseWriter, r *http.Request) {
	user, ok := socialUser(w, r)
	if !ok {
		return
	}
	items, err := h.storyService.ListNotifications(r.Context(), user, 60)
	if err != nil {
		RespondError(w, 500, "Không thể tải thông báo")
		return
	}
	RespondJSON(w, 200, map[string]any{"items": items})
}
func (h *StoryHandler) ReadNotification(w http.ResponseWriter, r *http.Request) {
	user, ok := socialUser(w, r)
	if !ok {
		return
	}
	if err := h.storyService.MarkNotificationRead(r.Context(), user, chi.URLParam(r, "id")); err != nil {
		RespondError(w, 500, "Không thể cập nhật thông báo")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *StoryHandler) Comments(w http.ResponseWriter, r *http.Request) {
	user, ok := socialUser(w, r)
	if !ok {
		return
	}
	storyID := chi.URLParam(r, "id")
	story, err := h.storyService.GetStory(r.Context(), storyID)
	if err != nil {
		RespondError(w, 404, "Không tìm thấy truyện")
		return
	}
	auth, _ := middleware.GetUserFromContext(r.Context())
	if story.Status != "published" && story.Status != "completed" && story.UserID != user && !auth.Role.IsStaff() {
		RespondError(w, 404, "Không tìm thấy truyện")
		return
	}
	items, err := h.storyService.ListComments(r.Context(), storyID, user)
	if err != nil {
		RespondError(w, 500, "Không thể tải bình luận")
		return
	}
	RespondJSON(w, 200, map[string]any{"items": items})
}
func (h *StoryHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, ok := socialUser(w, r)
	if !ok {
		return
	}
	var payload struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&payload); err != nil {
		RespondError(w, 400, "Bình luận không hợp lệ")
		return
	}
	body := strings.TrimSpace(payload.Body)
	if utf8.RuneCountInString(body) < 3 || utf8.RuneCountInString(body) > 2000 {
		RespondError(w, 400, "Bình luận phải có từ 3 đến 2000 ký tự")
		return
	}
	comment, err := h.storyService.CreateComment(r.Context(), chi.URLParam(r, "id"), user, body)
	if errors.Is(err, repository.ErrCommentNotFound) {
		RespondError(w, 404, "Truyện không tồn tại hoặc đã vượt giới hạn bình luận")
		return
	}
	if err != nil {
		RespondError(w, 500, "Không thể gửi bình luận")
		return
	}
	RespondJSON(w, http.StatusCreated, comment)
}
func (h *StoryHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	user, ok := socialUser(w, r)
	if !ok {
		return
	}
	err := h.storyService.DeleteComment(r.Context(), user, chi.URLParam(r, "id"))
	if errors.Is(err, repository.ErrCommentNotFound) {
		RespondError(w, 404, "Không tìm thấy bình luận của bạn")
		return
	}
	if err != nil {
		RespondError(w, 500, "Không thể xóa bình luận")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *StoryHandler) PendingComments(w http.ResponseWriter, r *http.Request) {
	if _, ok := socialUser(w, r); !ok {
		return
	}
	items, err := h.storyService.ListPendingComments(r.Context(), 80)
	if err != nil {
		RespondError(w, 500, "Không thể tải bình luận chờ duyệt")
		return
	}
	RespondJSON(w, 200, map[string]any{"items": items})
}
func (h *StoryHandler) ModerateComment(w http.ResponseWriter, r *http.Request) {
	if _, ok := socialUser(w, r); !ok {
		return
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 512)).Decode(&payload); err != nil ||
		(payload.Status != "approved" && payload.Status != "rejected") {
		RespondError(w, 400, "Trạng thái duyệt không hợp lệ")
		return
	}
	err := h.storyService.ModerateComment(r.Context(), chi.URLParam(r, "id"), payload.Status)
	if errors.Is(err, repository.ErrCommentNotFound) {
		RespondError(w, 404, "Không có bình luận chờ duyệt")
		return
	}
	if err != nil {
		RespondError(w, 500, "Không thể kiểm duyệt")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
