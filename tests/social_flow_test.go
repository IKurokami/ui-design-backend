package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ui-design-backend/internal/domain"
)

func TestSocialFollowNotificationsAndModeration(t *testing.T) {
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
		if cookie != nil {
			req.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	register := func(email string) *http.Cookie {
		w := request("POST", "/api/auth/register", map[string]any{"displayName": "Độc giả", "email": email, "password": "StrongPass123!", "confirmPassword": "StrongPass123!", "acceptedTerms": true}, nil)
		if w.Code != 200 {
			t.Fatalf("register failed %d: %s", w.Code, w.Body.String())
		}
		for _, c := range w.Result().Cookies() {
			if c.Name == "access_token" {
				return c
			}
		}
		t.Fatal("cookie missing")
		return nil
	}
	author := register("social-author@example.com")
	reader := register("social-reader@example.com")
	other := register("social-other@example.com")
	w := request("POST", "/api/stories", domain.CreateStoryDTO{Title: "Truyện theo dõi", Status: "published"}, author)
	if w.Code != 201 {
		t.Fatal(w.Body.String())
	}
	var story domain.Story
	if err := json.Unmarshal(w.Body.Bytes(), &story); err != nil {
		t.Fatal(err)
	}
	add := func(title string) domain.Chapter {
		w := request("POST", "/api/stories/"+story.ID+"/chapters", domain.CreateChapterDTO{Title: title, Content: "<p>Nội dung mới</p>"}, author)
		if w.Code != 201 {
			t.Fatalf("create chapter %d: %s", w.Code, w.Body.String())
		}
		var ch domain.Chapter
		if err := json.Unmarshal(w.Body.Bytes(), &ch); err != nil {
			t.Fatal(err)
		}
		return ch
	}
	add("Một")
	add("Hai")
	for _, path := range []string{"/api/follows/ids", "/api/notifications", "/api/stories/" + story.ID + "/comments"} {
		if got := request("GET", path, nil, nil).Code; got != 401 {
			t.Errorf("anonymous GET %s: %d", path, got)
		}
	}
	if got := request("POST", "/api/follows/"+story.ID, nil, nil).Code; got != 401 {
		t.Fatal(got)
	}
	for i := 0; i < 2; i++ {
		if got := request("POST", "/api/follows/"+story.ID, nil, reader).Code; got != 204 {
			t.Fatalf("follow %d", got)
		}
	}
	w = request("GET", "/api/follows/ids", nil, reader)
	var ids struct {
		IDs []string `json:"ids"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &ids); err != nil {
		t.Fatal(err)
	}
	if len(ids.IDs) != 1 || ids.IDs[0] != story.ID {
		t.Fatalf("follow ids = %+v", ids.IDs)
	}
	chapter := add("Ba")
	w = request("GET", "/api/notifications", nil, reader)
	var inbox struct {
		Items []domain.StoryNotification `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &inbox); err != nil {
		t.Fatal(err)
	}
	if len(inbox.Items) != 1 || inbox.Items[0].ChapterID != chapter.ID || inbox.Items[0].Read {
		t.Fatalf("inbox=%+v", inbox)
	}
	w = request("GET", "/api/notifications", nil, other)
	if err := json.Unmarshal(w.Body.Bytes(), &inbox); err != nil {
		t.Fatal(err)
	}
	if len(inbox.Items) != 0 {
		t.Fatal("notification leaked to other account")
	}
	if got := request("PUT", "/api/notifications/"+chapter.ID+"/read", nil, reader).Code; got != 204 {
		t.Fatal(got)
	}
	w = request("GET", "/api/notifications", nil, reader)
	if err := json.Unmarshal(w.Body.Bytes(), &inbox); err != nil {
		t.Fatal(err)
	}
	if len(inbox.Items) != 1 || !inbox.Items[0].Read {
		t.Fatalf("mark read failed %+v", inbox)
	}
	if got := request("DELETE", "/api/follows/"+story.ID, nil, reader).Code; got != 204 {
		t.Fatal(got)
	}
	add("Bốn")
	w = request("GET", "/api/notifications", nil, reader)
	if err := json.Unmarshal(w.Body.Bytes(), &inbox); err != nil {
		t.Fatal(err)
	}
	if len(inbox.Items) != 1 {
		t.Fatal("unfollow must stop new notifications")
	}

	w = request("POST", "/api/stories/"+story.ID+"/comments", map[string]string{"body": "Truyện rất cuốn!"}, reader)
	if w.Code != 201 {
		t.Fatalf("create comment %d %s", w.Code, w.Body.String())
	}
	var comment domain.StoryComment
	if err := json.Unmarshal(w.Body.Bytes(), &comment); err != nil {
		t.Fatal(err)
	}
	if comment.Status != "pending" {
		t.Fatalf("comment must be pending: %+v", comment)
	}
	w = request("GET", "/api/stories/"+story.ID+"/comments", nil, other)
	var comments struct {
		Items []domain.StoryComment `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &comments); err != nil {
		t.Fatal(err)
	}
	if len(comments.Items) != 0 {
		t.Fatal("unapproved comment leaked")
	}
	w = request("GET", "/api/stories/"+story.ID+"/comments", nil, reader)
	if err := json.Unmarshal(w.Body.Bytes(), &comments); err != nil {
		t.Fatal(err)
	}
	if len(comments.Items) != 1 || comments.Items[0].Status != "pending" {
		t.Fatal("commenter must see pending comment")
	}
	if got := request("GET", "/api/admin/comments", nil, reader).Code; got != 403 {
		t.Fatalf("user accessed moderation %d", got)
	}
	if got := request("PUT", "/api/admin/comments/"+comment.ID, map[string]string{"status": "approved"}, reader).Code; got != 403 {
		t.Fatalf("user could approve comment %d", got)
	}
	if got := request("DELETE", "/api/comments/"+comment.ID, nil, other).Code; got != 404 {
		t.Fatal(got)
	}
	if got := request("DELETE", "/api/comments/"+comment.ID, nil, reader).Code; got != 204 {
		t.Fatal(got)
	}
	if got := request("POST", "/api/stories/"+story.ID+"/comments", map[string]string{"body": "  "}, reader).Code; got != 400 {
		t.Fatal(got)
	}
}
