package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
	"ui-design-backend/internal/service"
)

type EbookHandler struct{ service *service.EbookService }

func NewEbookHandler(s *service.EbookService) *EbookHandler { return &EbookHandler{service: s} }

func ebookError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrEbookNotFound):
		RespondError(w, http.StatusNotFound, "Không tìm thấy ebook")
	case errors.Is(err, repository.ErrEbookConflict):
		RespondError(w, http.StatusConflict, "Ebook đã thay đổi ở cửa sổ khác. Tải lại trước khi lưu tiếp.")
	case errors.Is(err, service.ErrInvalidEbook):
		RespondError(w, http.StatusBadRequest, "Thiết kế ebook không hợp lệ. Nhập tên sách và kiểm tra các trang.")
	default:
		RespondError(w, http.StatusInternalServerError, "Không thể hoàn tất thao tác ebook")
	}
}
func ebookUser(w http.ResponseWriter, r *http.Request) (string, bool) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return "", false
	}
	return user.ID, true
}
func decodeEbook(w http.ResponseWriter, r *http.Request) (*domain.SaveEbookDTO, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)
	decoder := json.NewDecoder(r.Body)
	var req domain.SaveEbookDTO
	if err := decoder.Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu ebook không hợp lệ hoặc quá lớn")
		return nil, false
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		RespondError(w, http.StatusBadRequest, "Dữ liệu JSON không hợp lệ")
		return nil, false
	}
	return &req, true
}
func (h *EbookHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := ebookUser(w, r)
	if !ok {
		return
	}
	ebooks, err := h.service.List(r.Context(), userID)
	if err != nil {
		ebookError(w, err)
		return
	}
	RespondJSON(w, http.StatusOK, ebooks)
}
func (h *EbookHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := ebookUser(w, r)
	if !ok {
		return
	}
	ebook, err := h.service.Get(r.Context(), userID, chi.URLParam(r, "id"))
	if err != nil {
		ebookError(w, err)
		return
	}
	RespondJSON(w, http.StatusOK, ebook)
}
func (h *EbookHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := ebookUser(w, r)
	if !ok {
		return
	}
	req, ok := decodeEbook(w, r)
	if !ok {
		return
	}
	ebook, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		ebookError(w, err)
		return
	}
	RespondJSON(w, http.StatusCreated, ebook)
}
func (h *EbookHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := ebookUser(w, r)
	if !ok {
		return
	}
	req, ok := decodeEbook(w, r)
	if !ok {
		return
	}
	ebook, err := h.service.Update(r.Context(), userID, chi.URLParam(r, "id"), req)
	if err != nil {
		ebookError(w, err)
		return
	}
	RespondJSON(w, http.StatusOK, ebook)
}
func (h *EbookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := ebookUser(w, r)
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), userID, chi.URLParam(r, "id")); err != nil {
		ebookError(w, err)
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}
