package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadListingAndHistoricalActiveFilesBlocked(t *testing.T) {
	router, uploadDir, cleanup := setupStoryTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(
		`{"displayName":"Reader","email":"old-uploads-reader@example.com","password":"StrongPass123!","confirmPassword":"StrongPass123!","acceptedTerms":true}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("register: %d %s", w.Code, w.Body.String())
	}
	var access *http.Cookie
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "access_token" {
			access = cookie
		}
	}
	if access == nil {
		t.Fatal("missing cookie")
	}
	if err := os.WriteFile(filepath.Join(uploadDir, "old.svg"), []byte(`<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uploadDir, "index.html"), []byte("<script>alert(1)</script>"), 0644); err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{
		"/uploads/",
		"/uploads/old.svg",
		"/uploads/index.html",
		"/uploads/.gitignore",
	} {
		t.Run(target, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, target, nil)
			request.AddCookie(access)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusNotFound {
				t.Fatalf("%s should never be served, got %d: %s", target, recorder.Code, recorder.Body.String())
			}
		})
	}
}
