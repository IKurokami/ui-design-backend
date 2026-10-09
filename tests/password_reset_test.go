package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"ui-design-backend/internal/auth"
	"ui-design-backend/internal/config"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/repository"
	"ui-design-backend/internal/service"
	transporthttp "ui-design-backend/internal/transport/http"
)

type resetFixture struct {
	db     *repository.DB
	repo   *repository.PasswordResetRepo
	router http.Handler
}

func newResetFixture(t *testing.T, expiresAt time.Time, used bool) resetFixture {
	t.Helper()
	db, err := repository.NewDB("sqlite", filepath.Join(t.TempDir(), "reset.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	db.SetMaxOpenConns(1)
	ctx := context.Background()
	now := time.Now()
	users := repository.NewUserRepo(db)
	if err := users.Create(ctx, &domain.User{ID: "user", DisplayName: "Test", Email: "test@example.com", PasswordHash: "original", CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	sessions := repository.NewSessionRepo(db)
	if err := sessions.Create(ctx, &domain.Session{ID: "session", UserID: "user", TokenHash: "session-hash", CreatedAt: now, ExpiresAt: now.Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	repo := repository.NewPasswordResetRepo(db)
	reset := &domain.PasswordReset{ID: "reset", UserID: "user", TokenHash: auth.HashToken("test-token"), CreatedAt: now, ExpiresAt: expiresAt}
	if used {
		reset.UsedAt = &now
	}
	if err := repo.Create(ctx, reset); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{}
	svc := service.NewAuthService(cfg, users, sessions, repo, nil, nil)
	handler := transporthttp.NewAuthHandler(svc, nil)
	return resetFixture{db, repo, transporthttp.NewRouter(cfg, handler, nil, nil, nil)}
}

func resetRequest(router http.Handler, path string, body any) *httptest.ResponseRecorder {
	data, _ := json.Marshal(body)
	request := httptest.NewRequest(http.MethodPost, "/api/auth/password-reset/"+path, bytes.NewReader(data))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestPasswordResetLinkSingleUse(t *testing.T) {
	f := newResetFixture(t, time.Now().Add(time.Hour), false)
	for range 2 {
		response := resetRequest(f.router, "validate", map[string]string{"token": "test-token"})
		if response.Code != http.StatusNoContent {
			t.Fatalf("valid token: %d %s", response.Code, response.Body.String())
		}
		if response.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("validation must not be cached")
		}
	}
	response := resetRequest(f.router, "confirm", map[string]string{"token": "test-token", "newPassword": "NewPassword123"})
	if response.Code != http.StatusOK {
		t.Fatalf("first reset: %d %s", response.Code, response.Body.String())
	}
	if got := resetRequest(f.router, "validate", map[string]string{"token": "test-token"}); got.Code != http.StatusBadRequest {
		t.Fatalf("used link allowed: %d", got.Code)
	}
	if got := resetRequest(f.router, "confirm", map[string]string{"token": "test-token", "newPassword": "AnotherPassword123"}); got.Code != http.StatusBadRequest {
		t.Fatalf("token reuse allowed: %d", got.Code)
	}
	user, err := repository.NewUserRepo(f.db).FindByID(context.Background(), "user")
	if err != nil {
		t.Fatal(err)
	}
	if !auth.CheckPassword("NewPassword123", user.PasswordHash) {
		t.Fatal("password changed on token reuse")
	}
	session, err := repository.NewSessionRepo(f.db).FindByID(context.Background(), "session")
	if err != nil {
		t.Fatal(err)
	}
	if session.RevokedAt == nil {
		t.Fatal("old session was not revoked")
	}
}

func TestPasswordResetRejectsInvalidLinks(t *testing.T) {
	for _, tc := range []struct {
		name, token string
		expiresAt   time.Time
		used        bool
	}{
		{"missing", "", time.Now().Add(time.Hour), false},
		{"unknown", "unknown", time.Now().Add(time.Hour), false},
		{"expired", "test-token", time.Now().Add(-time.Hour), false},
		{"used", "test-token", time.Now().Add(time.Hour), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newResetFixture(t, tc.expiresAt, tc.used)
			for _, path := range []string{"validate", "confirm"} {
				response := resetRequest(f.router, path, map[string]string{"token": tc.token, "newPassword": "NewPassword123"})
				if response.Code != http.StatusBadRequest {
					t.Fatalf("%s: %d %s", path, response.Code, response.Body.String())
				}
			}
			user, err := repository.NewUserRepo(f.db).FindByID(context.Background(), "user")
			if err != nil {
				t.Fatal(err)
			}
			if user.PasswordHash != "original" {
				t.Fatal("invalid token changed password")
			}
		})
	}
}

func TestPasswordResetConcurrentConfirmation(t *testing.T) {
	f := newResetFixture(t, time.Now().Add(time.Hour), false)
	start := make(chan struct{})
	results := make(chan int, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			results <- resetRequest(f.router, "confirm", map[string]string{"token": "test-token", "newPassword": "NewPassword123"}).Code
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	counts := map[int]int{}
	for code := range results {
		counts[code]++
	}
	if counts[http.StatusOK] != 1 || counts[http.StatusBadRequest] != 1 {
		t.Fatalf("expected exactly one successful reset, got %v", counts)
	}
}

func TestPasswordResetCompletionRollback(t *testing.T) {
	for _, table := range []string{"users", "sessions"} {
		t.Run(table, func(t *testing.T) {
			f := newResetFixture(t, time.Now().Add(time.Hour), false)
			if _, err := f.db.Exec("CREATE TRIGGER fail_reset BEFORE UPDATE ON " + table + " BEGIN SELECT RAISE(ABORT, 'forced failure'); END"); err != nil {
				t.Fatal(err)
			}
			response := resetRequest(f.router, "confirm", map[string]string{"token": "test-token", "newPassword": "NewPassword123"})
			if response.Code != http.StatusInternalServerError {
				t.Fatalf("database failure should return 500: %d", response.Code)
			}
			reset, err := f.repo.FindByTokenHash(context.Background(), auth.HashToken("test-token"))
			if err != nil {
				t.Fatal(err)
			}
			if reset.UsedAt != nil {
				t.Fatal("failed reset consumed token")
			}
			user, err := repository.NewUserRepo(f.db).FindByID(context.Background(), "user")
			if err != nil {
				t.Fatal(err)
			}
			if user.PasswordHash != "original" {
				t.Fatal("failed reset changed password")
			}
			if _, err := f.db.Exec("DROP TRIGGER fail_reset"); err != nil {
				t.Fatal(err)
			}
			if err := f.repo.Complete(context.Background(), "reset", "retry-hash"); err != nil {
				t.Fatalf("retry failed: %v", err)
			}
		})
	}
}

func TestPasswordResetCompletionRechecksExpiry(t *testing.T) {
	f := newResetFixture(t, time.Now().Add(time.Hour), false)
	if response := resetRequest(f.router, "validate", map[string]string{"token": "test-token"}); response.Code != http.StatusNoContent {
		t.Fatal("initial validation failed")
	}
	if _, err := f.db.Exec("UPDATE password_resets SET expires_at = ? WHERE id = ?", time.Now().Add(-time.Second), "reset"); err != nil {
		t.Fatal(err)
	}
	if err := f.repo.Complete(context.Background(), "reset", "new-hash"); !errors.Is(err, repository.ErrResetNotFound) {
		t.Fatalf("expired token accepted: %v", err)
	}
	reset, err := f.repo.FindByTokenHash(context.Background(), auth.HashToken("test-token"))
	if err != nil {
		t.Fatal(err)
	}
	if reset.UsedAt != nil {
		t.Fatal("expired token mutation was not rolled back")
	}
}
