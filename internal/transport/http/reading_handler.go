package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
)

// A reader may never inspect or persist progress for another author's draft.
func (h *StoryHandler) mayReadProgress(w http.ResponseWriter, r *http.Request) bool {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok { RespondError(w, http.StatusUnauthorized, "Chưa đăng nhập"); return false }
	story, err := h.storyService.GetStory(r.Context(), chi.URLParam(r, "id"))
	if err != nil { RespondError(w, http.StatusNotFound, "Không tìm thấy truyện"); return false }
	if story.Status != "published" && story.Status != "completed" && story.UserID != user.ID && !user.Role.IsStaff() {
		RespondError(w, http.StatusNotFound, "Không tìm thấy truyện")
		return false
	}
	return true
}

func (h *StoryHandler) GetReadingProgress(w http.ResponseWriter, r *http.Request) {
	if !h.mayReadProgress(w, r) { return }
	w.Header().Set("Cache-Control", "private, no-store")
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}
	progress, err := h.storyService.GetReadingProgress(r.Context(), user.ID, chi.URLParam(r, "id"))
	if err != nil {
		h.readingError(w, err)
		return
	}
	RespondJSON(w, http.StatusOK, progress)
}

func (h *StoryHandler) SaveReadingProgress(w http.ResponseWriter, r *http.Request) {
	if !h.mayReadProgress(w, r) { return }
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}
	var progress domain.ReadingProgress
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&progress); err != nil || progress.ChapterID == "" {
		RespondError(w, http.StatusBadRequest, "Thiếu ID chương hợp lệ")
		return
	}
	if err := h.storyService.SaveReadingProgress(r.Context(), user.ID, chi.URLParam(r, "id"), progress.ChapterID); err != nil {
		h.readingError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *StoryHandler) readingError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrStoryNotFound) || errors.Is(err, repository.ErrChapterNotFound) {
		RespondError(w, http.StatusNotFound, "Không tìm thấy truyện hoặc chương có nội dung")
		return
	}
	RespondError(w, http.StatusInternalServerError, "Không thể lưu hoặc tải chương đọc gần nhất")
}
