package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"ui-design-backend/internal/auth"
	"ui-design-backend/internal/config"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
	"ui-design-backend/internal/service"
	transporthttp "ui-design-backend/internal/transport/http"
)

func TestManagementRoles(t *testing.T) {
	db, err := repository.NewDB("sqlite", filepath.Join(t.TempDir(), "admin.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	cfg := &config.Config{Env: "test", ClientOrigin: "http://localhost:5173", JWTSecret: "admin-test-signing-key", AccessTokenTTL: time.Hour, RefreshTokenRememberTTL: time.Hour}
	users := repository.NewUserRepo(db)
	admin := &domain.User{ID: uuid.NewString(), DisplayName: "Admin", Email: "admin@example.com", PasswordHash: "x", Role: domain.RoleAdmin, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	moderator := &domain.User{ID: uuid.NewString(), DisplayName: "Moderator", Email: "mod@example.com", PasswordHash: "x", Role: domain.RoleModerator, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	reader := &domain.User{ID: uuid.NewString(), DisplayName: "Reader", Email: "reader@example.com", PasswordHash: "x", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	for _, user := range []*domain.User{admin, moderator, reader} {
		if err := users.Create(context.Background(), user); err != nil {
			t.Fatal(err)
		}
	}
	stories := repository.NewStoryRepo(db)
	storyService := service.NewStoryService(stories, users)
	management := service.NewManagementService(users, storyService)
	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL)
	authMW := middleware.NewAuthMiddleware(cfg, tokens, nil, users)
	authHandler := transporthttp.NewAuthHandler(nil, authMW, transporthttp.NewAdminHandler(management))
	router := transporthttp.NewRouter(cfg, authHandler, nil, nil, authMW)

	request := func(method, path, token string) *httptest.ResponseRecorder {
		tokenValue, _, err := tokens.GenerateAccessToken(token, token+"@example.com", "User", "session")
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(method, path, nil)
		req.Header.Set("Authorization", "Bearer "+tokenValue)
		result := httptest.NewRecorder()
		router.ServeHTTP(result, req)
		return result
	}

	if response := request(http.MethodGet, "/api/admin/users", moderator.ID); response.Code != http.StatusForbidden {
		t.Fatalf("moderator users status = %d", response.Code)
	}
	if response := request(http.MethodGet, "/api/admin/stories", moderator.ID); response.Code != http.StatusOK {
		t.Fatalf("moderator stories status = %d", response.Code)
	}
	if response := request(http.MethodGet, "/api/admin/users", admin.ID); response.Code != http.StatusOK {
		t.Fatalf("admin users status = %d", response.Code)
	}
	if response := request(http.MethodDelete, "/api/admin/users/"+reader.ID, admin.ID); response.Code != http.StatusNoContent {
		t.Fatalf("admin delete user status = %d", response.Code)
	}
	if _, err := users.FindByID(context.Background(), reader.ID); err == nil {
		t.Fatal("deleted user still exists")
	}
}
