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

type StoryHandler struct {
	storyService *service.StoryService
}

func NewStoryHandler(storyService *service.StoryService) *StoryHandler {
	return &StoryHandler{
		storyService: storyService,
	}
}

func (h *StoryHandler) CreateStory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}

	var req domain.CreateStoryDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu JSON không hợp lệ")
		return
	}

	story, err := h.storyService.CreateStory(r.Context(), user.ID, &req)
	if err != nil {
		if errors.Is(err, service.ErrStoryTitleRequired) {
			RespondError(w, http.StatusBadRequest, "Tên truyện không được để trống")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi máy chủ khi tạo truyện: "+err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, story)
}

func (h *StoryHandler) ListStories(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("paginated") == "true" || r.Header.Get("X-MOC-Infinite") == "true" {
		page, size := parseStoryPage(r)
		readerID := ""
		ownerID := ""
		if r.URL.Query().Get("reading") == "true" {
			user, ok := middleware.GetUserFromContext(r.Context())
			if !ok {
				RespondError(w, 401, "Chưa xác thực")
				return
			}
			readerID = user.ID
			w.Header().Set("Cache-Control", "private, no-store")
		} else if r.URL.Query().Get("mine") == "true" {
			user, ok := middleware.GetUserFromContext(r.Context())
			if !ok {
				RespondError(w, 401, "Chưa xác thực")
				return
			}
			ownerID = user.ID
		}
		items, err := h.storyService.ListPage(r.Context(), ownerID, readerID, page, size+1)
		if err != nil {
			RespondError(w, 500, "Lỗi khi lấy danh sách truyện")
			return
		}
		more := len(items) > size
		if more {
			items = items[:size]
		}
		next := 0
		if more {
			next = page + 1
		}
		RespondJSON(w, 200, map[string]any{"items": items, "nextPage": next, "hasMore": more})
		return
	}
	mine := r.URL.Query().Get("mine")
	var userID string
	if mine == "true" || mine == "1" {
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
			return
		}
		userID = user.ID
	}

	var stories []*domain.Story
	var err error
	if r.URL.Query().Get("reading") == "true" {
		w.Header().Set("Cache-Control", "private, no-store")
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
			return
		}
		stories, err = h.storyService.ListReadingStories(r.Context(), user.ID)
	} else {
		stories, err = h.storyService.ListStories(r.Context(), userID)
	}
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Lỗi khi lấy danh sách truyện")
		return
	}
	if stories == nil {
		stories = []*domain.Story{}
	}

	RespondJSON(w, http.StatusOK, stories)
}

func parseStoryPage(r *http.Request) (int, int) {
	p, _ := strconv.Atoi(r.URL.Query().Get("page"))
	n, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if p < 1 {
		p = 1
	}
	if n < 1 {
		n = 24
	}
	if n > 60 {
		n = 60
	}
	return p, n
}

func (h *StoryHandler) GetMyStories(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}

	if r.URL.Query().Get("paginated") == "true" || r.Header.Get("X-MOC-Infinite") == "true" {
		page, size := parseStoryPage(r)
		items, err := h.storyService.ListPage(r.Context(), user.ID, "", page, size+1)
		if err != nil {
			RespondError(w, 500, "Lỗi khi lấy danh sách truyện trong tủ sách")
			return
		}
		more := len(items) > size
		if more {
			items = items[:size]
		}
		next := 0
		if more {
			next = page + 1
		}
		RespondJSON(w, 200, map[string]any{"items": items, "nextPage": next, "hasMore": more})
		return
	}
	stories, err := h.storyService.ListStories(r.Context(), user.ID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Lỗi khi lấy danh sách truyện trong tủ sách")
		return
	}
	if stories == nil {
		stories = []*domain.Story{}
	}

	RespondJSON(w, http.StatusOK, stories)
}

func (h *StoryHandler) GetStory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/stories/")
		if idx := strings.Index(id, "/"); idx != -1 {
			id = id[:idx]
		}
	}
	if id == "" {
		RespondError(w, http.StatusBadRequest, "Thiếu ID truyện")
		return
	}

	story, err := h.storyService.GetStory(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrStoryNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy truyện")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi máy chủ khi lấy thông tin truyện")
		return
	}

	RespondJSON(w, http.StatusOK, story)
}

func (h *StoryHandler) UpdateStory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/stories/")
	}
	if id == "" {
		RespondError(w, http.StatusBadRequest, "Thiếu ID truyện")
		return
	}

	var req domain.UpdateStoryDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu JSON không hợp lệ")
		return
	}

	story, err := h.storyService.UpdateStory(r.Context(), user.ID, id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrStoryNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy truyện")
			return
		}
		if err.Error() == "forbidden" {
			RespondError(w, http.StatusForbidden, "Bạn không có quyền chỉnh sửa truyện này")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi máy chủ khi cập nhật truyện")
		return
	}

	RespondJSON(w, http.StatusOK, story)
}

func (h *StoryHandler) DeleteStory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/stories/")
	}
	if id == "" {
		RespondError(w, http.StatusBadRequest, "Thiếu ID truyện")
		return
	}

	if err := h.storyService.DeleteStory(r.Context(), user.ID, id); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy truyện")
			return
		}
		if errors.Is(err, repository.ErrStoryNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy truyện")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi khi xóa truyện")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StoryHandler) CreateChapter(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}

	storyID := chi.URLParam(r, "id")
	if storyID == "" {
		path := strings.TrimPrefix(r.URL.Path, "/api/stories/")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 {
			storyID = parts[0]
		}
	}
	if storyID == "" {
		RespondError(w, http.StatusBadRequest, "Đường dẫn không hợp lệ")
		return
	}

	var req domain.CreateChapterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu JSON không hợp lệ")
		return
	}

	chapter, err := h.storyService.CreateChapter(r.Context(), user.ID, storyID, &req)
	if err != nil {
		if errors.Is(err, repository.ErrStoryNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy truyện")
			return
		}
		if errors.Is(err, service.ErrChapterTitleRequired) {
			RespondError(w, http.StatusBadRequest, "Tiêu đề chương không được để trống")
			return
		}
		if err.Error() == "forbidden" {
			RespondError(w, http.StatusForbidden, "Bạn không có quyền thêm chương cho truyện này")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi khi tạo chương")
		return
	}

	RespondJSON(w, http.StatusCreated, chapter)
}

func (h *StoryHandler) GetChapter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/chapters/")
	}
	if id == "" {
		RespondError(w, http.StatusBadRequest, "Thiếu ID chương")
		return
	}

	chapter, err := h.storyService.GetChapter(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrChapterNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy chương")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi khi lấy thông tin chương")
		return
	}

	RespondJSON(w, http.StatusOK, chapter)
}

func (h *StoryHandler) UpdateChapter(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/chapters/")
	}
	if id == "" {
		RespondError(w, http.StatusBadRequest, "Thiếu ID chương")
		return
	}

	var req domain.UpdateChapterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu JSON không hợp lệ")
		return
	}

	chapter, err := h.storyService.UpdateChapter(r.Context(), user.ID, id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrChapterNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy chương")
			return
		}
		if err.Error() == "forbidden" {
			RespondError(w, http.StatusForbidden, "Bạn không có quyền cập nhật chương này")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi khi cập nhật chương")
		return
	}

	RespondJSON(w, http.StatusOK, chapter)
}

func (h *StoryHandler) DeleteChapter(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/chapters/")
	}
	if id == "" {
		RespondError(w, http.StatusBadRequest, "Thiếu ID chương")
		return
	}

	if err := h.storyService.DeleteChapter(r.Context(), user.ID, id); err != nil {
		if errors.Is(err, repository.ErrChapterNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy chương")
			return
		}
		if err.Error() == "forbidden" {
			RespondError(w, http.StatusForbidden, "Bạn không có quyền xóa chương này")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi khi xóa chương")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
