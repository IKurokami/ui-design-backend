package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
	"ui-design-backend/internal/service"
)

type ProfileHandler struct {
	profileService *service.ProfileService
	UploadDir      string
}

func NewProfileHandler(profileService *service.ProfileService, uploadDir string) *ProfileHandler {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	_ = os.MkdirAll(uploadDir, 0755)
	return &ProfileHandler{
		profileService: profileService,
		UploadDir:      uploadDir,
	}
}

func (h *ProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}

	profile, err := h.profileService.GetOrCreateProfile(r.Context(), user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy người dùng")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi máy chủ khi lấy profile")
		return
	}

	RespondJSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Chưa xác thực")
		return
	}

	var dto domain.UpdateProfileDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		RespondError(w, http.StatusBadRequest, "Dữ liệu JSON không hợp lệ")
		return
	}

	profile, err := h.profileService.UpdateProfile(r.Context(), user.ID, &dto)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy người dùng")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi máy chủ khi cập nhật profile")
		return
	}

	RespondJSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	if userID == "" {
		RespondError(w, http.StatusBadRequest, "Thiếu tham số userId")
		return
	}

	profile, err := h.profileService.GetPublicProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			RespondError(w, http.StatusNotFound, "Không tìm thấy hồ sơ người dùng")
			return
		}
		RespondError(w, http.StatusInternalServerError, "Lỗi máy chủ khi lấy hồ sơ người dùng")
		return
	}

	RespondJSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Giới hạn dung lượng tối đa 10 MB
	const maxUploadSize = 10 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		RespondError(w, http.StatusBadRequest, "Kích thước file quá lớn (tối đa 10MB) hoặc dữ liệu gửi lên không đúng định dạng multipart")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		file, header, err = r.FormFile("avatar")
	}
	if err != nil {
		file, header, err = r.FormFile("banner")
	}
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Không tìm thấy file tải lên (hỗ trợ trường 'file', 'avatar' hoặc 'banner')")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
		".svg":  true,
	}

	if !allowedExts[ext] {
		RespondError(w, http.StatusBadRequest, "Định dạng file không được hỗ trợ. Chỉ chấp nhận các định dạng: JPG, PNG, WEBP, GIF, SVG.")
		return
	}

	if err := os.MkdirAll(h.UploadDir, 0755); err != nil {
		RespondError(w, http.StatusInternalServerError, "Không thể khởi tạo thư mục lưu file")
		return
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	targetFilePath := filepath.Join(h.UploadDir, filename)

	dst, err := os.Create(targetFilePath)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Không thể lưu file trên máy chủ")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		RespondError(w, http.StatusInternalServerError, "Lỗi xảy ra khi ghi file")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]any{
		"url":      "/uploads/" + filename,
		"filename": filename,
		"size":     header.Size,
	})
}
