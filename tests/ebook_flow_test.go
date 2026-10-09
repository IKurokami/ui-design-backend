package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"ui-design-backend/internal/auth"
	"ui-design-backend/internal/config"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
	"ui-design-backend/internal/service"
	transporthttp "ui-design-backend/internal/transport/http"
)

func setupEbookTestApp(t *testing.T) (http.Handler, *http.Cookie) {
	t.Helper()
	db, err := repository.NewDB("sqlite", filepath.Join(t.TempDir(), "ebooks.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	cfg := &config.Config{Env: "test", ClientOrigin: "http://localhost:5173", JWTSecret: "ebook-test-signing-key", AccessTokenTTL: time.Hour}
	now := time.Now().UTC()
	user := &domain.User{ID: "ebook-test-user", DisplayName: "Ebook test", Email: "ebook-test@example.test", PasswordHash: "test-only", CreatedAt: now, UpdatedAt: now}
	if err := repository.NewUserRepo(db).Create(context.Background(), user); err != nil {
		t.Fatal(err)
	}
	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL)
	token, _, err := tokens.GenerateAccessToken(user.ID, user.Email, user.DisplayName, "ebook-test-session")
	if err != nil {
		t.Fatal(err)
	}
	handler := transporthttp.NewEbookHandler(service.NewEbookService(repository.NewEbookRepo(db)))
	router := transporthttp.NewRouter(cfg, nil, nil, nil, middleware.NewAuthMiddleware(cfg, tokens, nil, repository.NewUserRepo(db)), handler)
	return router, &http.Cookie{Name: "access_token", Value: token}
}

func ebookRequest(t *testing.T, router http.Handler, cookie *http.Cookie, method, path string, payload any) *httptest.ResponseRecorder {
	t.Helper()
	var body []byte
	if payload != nil {
		var err error
		body, err = json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestEbookCanvasSaveAndReload(t *testing.T) {
	router, cookie := setupEbookTestApp(t)
	design := json.RawMessage(`{"version":2,"title":"Sách của tôi","pageBackground":"#f8f2e8","accent":"#435a41","pages":[{"id":"cover","kind":"cover","background":"#f8f2e8","accent":"#435a41","elements":[{"id":"title","kind":"text","html":"Sách của tôi","x":130,"y":1300,"fontSize":88,"color":"#254c68","rotation":15}]},{"id":"page-1","kind":"blank","background":"#f8f2e8","accent":"#435a41","elements":[]}]}`)
	payload := domain.SaveEbookDTO{Title: "Sách của tôi", Design: design}
	w := ebookRequest(t, router, cookie, http.MethodPost, "/api/ebooks", payload)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", w.Code, w.Body.String())
	}
	var created domain.Ebook
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Revision != 1 || created.PageCount != 2 {
		t.Fatalf("unexpected created ebook: %+v", created)
	}
	payload.Title = "Sách đã chỉnh sửa"
	payload.ExpectedRevision = created.Revision
	w = ebookRequest(t, router, cookie, http.MethodPut, "/api/ebooks/"+created.ID, payload)
	if w.Code != http.StatusOK {
		t.Fatalf("update: %d %s", w.Code, w.Body.String())
	}
	w = ebookRequest(t, router, cookie, http.MethodGet, "/api/ebooks/"+created.ID, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("reload: %d %s", w.Code, w.Body.String())
	}
	var reloaded domain.Ebook
	if err := json.Unmarshal(w.Body.Bytes(), &reloaded); err != nil {
		t.Fatal(err)
	}
	if reloaded.Title != payload.Title || reloaded.Revision != 2 || !bytes.Equal(reloaded.Design, design) {
		t.Fatalf("saved canvas changed on reload: %+v", reloaded)
	}
	w = ebookRequest(t, router, cookie, http.MethodGet, "/api/ebooks", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("library: %d %s", w.Code, w.Body.String())
	}
	var library []domain.Ebook
	if err := json.Unmarshal(w.Body.Bytes(), &library); err != nil || len(library) != 1 || library[0].ID != created.ID {
		t.Fatalf("ebook missing from library: %s", w.Body.String())
	}
	w = ebookRequest(t, router, cookie, http.MethodPut, "/api/ebooks/"+created.ID, payload)
	if w.Code != http.StatusConflict {
		t.Fatalf("stale save should conflict: %d %s", w.Code, w.Body.String())
	}
}
