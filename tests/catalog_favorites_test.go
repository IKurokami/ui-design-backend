package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"ui-design-backend/internal/domain"
)

func TestCatalogSearchFavoritesAndDraftProtection(t *testing.T) {
	router, _, cleanup := setupStoryTestServer(t)
	defer cleanup()

	request := func(method, path string, payload any, cookie *http.Cookie) *httptest.ResponseRecorder {
		var body bytes.Buffer
		if payload != nil {
			if err := json.NewEncoder(&body).Encode(payload); err != nil {
				t.Fatal(err)
			}
		}
		req := httptest.NewRequest(method, path, &body)
		req.Header.Set("Content-Type", "application/json")
		if cookie != nil {
			req.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	register := func(email string) *http.Cookie {
		w := request("POST", "/api/auth/register", map[string]any{
			"displayName": "Độc giả", "email": email, "password": "StrongPass123!",
			"confirmPassword": "StrongPass123!", "acceptedTerms": true,
		}, nil)
		if w.Code != 200 {
			t.Fatalf("register failed: %d %s", w.Code, w.Body.String())
		}
		for _, cookie := range w.Result().Cookies() {
			if cookie.Name == "access_token" {
				return cookie
			}
		}
		t.Fatal("Missing cookie")
		return nil
	}
	owner := register("catalog-owner@example.com")
	reader := register("catalog-reader@example.com")

	makeStory := func(title, category, tags, status string) (domain.Story, domain.Chapter) {
		w := request("POST", "/api/stories", domain.CreateStoryDTO{
			Title: title, Category: category, Tags: tags, Status: status,
		}, owner)
		if w.Code != 201 {
			t.Fatalf("create story: %d %s", w.Code, w.Body.String())
		}
		var story domain.Story
		if err := json.Unmarshal(w.Body.Bytes(), &story); err != nil {
			t.Fatal(err)
		}
		var chapter domain.Chapter
		for i := 1; i <= 2; i++ {
			w = request("POST", "/api/stories/"+story.ID+"/chapters", domain.CreateChapterDTO{
				ChapterNumber: i, Title: "Chương", Content: "<p>Nội dung.</p>",
			}, owner)
			if w.Code != 201 {
				t.Fatalf("chapter: %d %s", w.Code, w.Body.String())
			}
			if err := json.Unmarshal(w.Body.Bytes(), &chapter); err != nil {
				t.Fatal(err)
			}
		}
		return story, chapter
	}
	first, chapter := makeStory("Bí mật mặt trăng", "Kỳ Ảo", "Phiêu Lưu, Huyền Ảo", "published")
	second, _ := makeStory("Bình minh", "Đời Thường", "Bách Hợp", "completed")
	draft, secretChapter := makeStory("Bản thảo bí mật", "Kỳ Ảo", "Bí mật", "draft")

	for _, path := range []string{"/api/stories", "/api/stories?paginated=true&q=bi", "/api/stories/" + first.ID, "/api/chapters/" + chapter.ID, "/api/favorites/ids"} {
		if got := request("GET", path, nil, nil).Code; got != 401 {
			t.Errorf("unauthenticated GET %s = %d, want 401", path, got)
		}
	}
	for _, path := range []string{"/api/stories/" + draft.ID, "/api/chapters/" + secretChapter.ID} {
		if got := request("GET", path, nil, reader).Code; got != 404 {
			t.Errorf("other reader GET draft %s = %d, want 404", path, got)
		}
		if got := request("GET", path, nil, owner).Code; got != 200 {
			t.Errorf("owner GET draft %s = %d, want 200", path, got)
		}
	}

	search := func(query string) []domain.Story {
		w := request("GET", "/api/stories?paginated=true&page=1&pageSize=1&q="+url.QueryEscape(query), nil, reader)
		if w.Code != 200 {
			t.Fatalf("search: %d %s", w.Code, w.Body.String())
		}
		var p struct {
			Items   []domain.Story `json:"items"`
			HasMore bool           `json:"hasMore"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &p); err != nil {
			t.Fatal(err)
		}
		return p.Items
	}
	for _, query := range []string{"mặt trăng", "Phiêu Lưu", "Kỳ Ảo"} {
		got := search(query)
		if len(got) != 1 || got[0].ID != first.ID {
			t.Errorf("query %q returned %+v", query, got)
		}
	}
	if got := search("Bách Hợp"); len(got) != 1 || got[0].ID != second.ID {
		t.Fatalf("tag search failed: %+v", got)
	}
	if got := search("Bản thảo"); len(got) != 0 {
		t.Fatalf("draft leaked into search: %+v", got)
	}

	favorite := "/api/favorites/" + first.ID
	if got := request("POST", favorite, nil, nil).Code; got != 401 {
		t.Fatalf("unauth favorite: %d", got)
	}
	if got := request("POST", "/api/favorites/"+draft.ID, nil, reader).Code; got != 404 {
		t.Fatalf("draft favorited: %d", got)
	}
	for i := 0; i < 2; i++ {
		if got := request("POST", favorite, nil, reader).Code; got != 204 {
			t.Fatalf("favorite: %d", got)
		}
	}
	if got := request("POST", "/api/favorites/"+second.ID, nil, owner).Code; got != 204 {
		t.Fatal(got)
	}

	w := request("GET", "/api/favorites/ids", nil, reader)
	var ids struct {
		IDs []string `json:"ids"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &ids); err != nil {
		t.Fatal(err)
	}
	if len(ids.IDs) != 1 || ids.IDs[0] != first.ID {
		t.Fatalf("favorite ids leaked accounts: %+v", ids.IDs)
	}
	w = request("GET", "/api/favorites?page=1&pageSize=24", nil, reader)
	var list struct {
		Items   []domain.Story `json:"items"`
		HasMore bool           `json:"hasMore"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].ID != first.ID || list.HasMore {
		t.Fatalf("favorite list: %+v", list)
	}
	if got := request("DELETE", favorite, nil, reader).Code; got != 204 {
		t.Fatal(got)
	}
	if got := request("DELETE", favorite, nil, reader).Code; got != 204 {
		t.Fatal(got)
	}
	w = request("GET", "/api/favorites/ids", nil, reader)
	if err := json.Unmarshal(w.Body.Bytes(), &ids); err != nil {
		t.Fatal(err)
	}
	if len(ids.IDs) != 0 {
		t.Fatalf("favorite still present: %+v", ids.IDs)
	}
}
