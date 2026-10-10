package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/middleware"
)

func (h *StoryHandler) favoriteUser(w http.ResponseWriter, r *http.Request) (string, bool) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa đăng nhập")
		return "", false
	}
	w.Header().Set("Cache-Control", "private, no-store")
	return user.ID, true
}

func (h *StoryHandler) GetFavoriteIDs(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.favoriteUser(w, r)
	if !ok { return }
	ids, err := h.storyService.FavoriteIDs(r.Context(), userID)
	if err != nil { RespondError(w, http.StatusInternalServerError, "Không thể tải truyện đã lưu"); return }
	RespondJSON(w, http.StatusOK, map[string]any{"ids": ids})
}

func (h *StoryHandler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.favoriteUser(w, r)
	if !ok { return }
	page, size := parseStoryPage(r)
	items, err := h.storyService.ListFavorites(r.Context(), userID, page, size+1)
	if err != nil { RespondError(w, http.StatusInternalServerError, "Không thể tải truyện đã lưu"); return }
	if items == nil { items = []*domain.Story{} }
	more := len(items) > size
	if more { items = items[:size] }
	next := 0
	if more { next = page+1 }
	RespondJSON(w, http.StatusOK, map[string]any{"items": items, "hasMore": more, "nextPage": next})
}

func (h *StoryHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.favoriteUser(w, r)
	if !ok { return }
	id := chi.URLParam(r, "id")
	if id == "" { RespondError(w, http.StatusBadRequest, "Thiếu ID truyện"); return }
	saved, err := h.storyService.SaveFavorite(r.Context(), userID, id)
	if err != nil { RespondError(w, http.StatusInternalServerError, "Không thể lưu truyện"); return }
	if !saved { RespondError(w, http.StatusNotFound, "Không tìm thấy truyện có thể lưu"); return }
	w.WriteHeader(http.StatusNoContent)
}

func (h *StoryHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.favoriteUser(w, r)
	if !ok { return }
	id := chi.URLParam(r, "id")
	if id == "" { RespondError(w, http.StatusBadRequest, "Thiếu ID truyện"); return }
	if err := h.storyService.RemoveFavorite(r.Context(), userID, id); err != nil {
		RespondError(w, http.StatusInternalServerError, "Không thể bỏ lưu truyện")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
