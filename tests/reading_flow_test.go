package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"ui-design-backend/internal/domain"
)

func TestReadingProgressAndPublicLibrary(t *testing.T) {
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
		w := request("POST", "/api/auth/register", map[string]any{"displayName": "Người đọc", "email": email, "password": "StrongPass123!", "confirmPassword": "StrongPass123!", "acceptedTerms": true}, nil)
		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}
		for _, cookie := range w.Result().Cookies() {
			if cookie.Name == "access_token" {
				return cookie
			}
		}
		t.Fatal("Missing auth cookie")
		return nil
	}
	author, reader := register("reading-author@example.com"), register("reading-reader@example.com")
	createStory := func(title, status string) domain.Story {
		w := request("POST", "/api/stories", domain.CreateStoryDTO{Title: title, Status: status}, author)
		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}
		var story domain.Story
		if err := json.Unmarshal(w.Body.Bytes(), &story); err != nil {
			t.Fatal(err)
		}
		return story
	}
	addChapter := func(story domain.Story, number int, content string) domain.Chapter {
		w := request("POST", "/api/stories/"+story.ID+"/chapters", domain.CreateChapterDTO{ChapterNumber: number, Title: "Chương", Content: content}, author)
		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}
		var ch domain.Chapter
		if err := json.Unmarshal(w.Body.Bytes(), &ch); err != nil {
			t.Fatal(err)
		}
		return ch
	}
	list := func() []domain.Story {
		w := request("GET", "/api/stories", nil, reader)
		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}
		var stories []domain.Story
		if err := json.Unmarshal(w.Body.Bytes(), &stories); err != nil {
			t.Fatal(err)
		}
		return stories
	}
	story := createStory("Truyện có nội dung", "published")
	first := addChapter(story, 1, "<p>Nội dung chương 1.</p>")
	second := addChapter(story, 2, "<p><br></p><p>&nbsp;\u200b</p>")
	if len(list()) != 0 {
		t.Fatal("Empty rich-text chapter must not make a story eligible")
	}
	w := request("PUT", "/api/chapters/"+second.ID, map[string]string{"content": "<p>Nội dung chương 2.</p>"}, author)
	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
	if len(list()) != 1 {
		t.Fatal("Two readable chapters should appear")
	}
	draft := createStory("Bản nháp", "draft")
	addChapter(draft, 1, "Một")
	other := addChapter(draft, 2, "Hai")
	if len(list()) != 1 {
		t.Fatal("Draft must stay out of public library")
	}
	w = request("GET", "/api/stories/mine", nil, author)
	var mine []domain.Story
	_ = json.Unmarshal(w.Body.Bytes(), &mine)
	if len(mine) != 2 {
		t.Fatal("Author bookshelf must retain draft and public story")
	}
	path := "/api/stories/" + story.ID + "/reading-progress"
	getProgress := func(cookie *http.Cookie) domain.ReadingProgress {
		w := request("GET", path, nil, cookie)
		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}
		var p domain.ReadingProgress
		if err := json.Unmarshal(w.Body.Bytes(), &p); err != nil {
			t.Fatal(err)
		}
		return p
	}
	if getProgress(reader).ChapterID != "" {
		t.Fatal("New reader should have no progress")
	}
	for _, id := range []string{first.ID, second.ID, first.ID} {
		w = request("PUT", path, map[string]string{"chapterId": id}, reader)
		if w.Code != 204 {
			t.Fatal(w.Body.String())
		}
		if getProgress(reader).ChapterID != id {
			t.Fatal("Must persist the last read chapter, including reading backwards")
		}
	}

	if getProgress(author).ChapterID != "" {
		t.Fatal("Reading history must be isolated per user")
	}
	for _, c := range []*http.Cookie{reader, author, nil} {
		w := request("GET", "/api/stories?reading=true", nil, c)
		if c == nil {
			if w.Code != 401 {
				t.Fatal("Reading list must require authentication")
			}
			continue
		}
		var reading []domain.Story
		if err := json.Unmarshal(w.Body.Bytes(), &reading); err != nil {
			t.Fatal(err)
		}
		expected := 0
		if c == reader {
			expected = 1
		}
		if w.Code != 200 || len(reading) != expected {
			t.Fatal("Reading list must use this account's history")
		}
		if c == reader {
			progress := reading[0].ReadingProgress
			if progress == nil || progress.ChapterID != first.ID || progress.ChapterNumber != 1 || progress.ChapterTitle != first.Title || progress.UpdatedAt == nil || time.Since(*progress.UpdatedAt) > time.Minute {
				t.Fatalf("Reading API must expose the saved chapter and server timestamp: %+v", progress)
			}
		}
	}
	for _, method := range []string{"GET", "PUT"} {
		if request(method, path, map[string]string{"chapterId": first.ID}, nil).Code != 401 {
			t.Fatal("Progress must require authentication")
		}
	}
	if request("PUT", path, map[string]string{"chapterId": other.ID}, reader).Code != 404 {
		t.Fatal("Cannot save a chapter from another story")
	}
	if request("PUT", path, map[string]string{"chapterId": ""}, reader).Code != 400 {
		t.Fatal("Empty chapter ID must be rejected")
	}
	// Unauthorized story deletion must not destroy chapters or another reader's progress.
	if request("DELETE", "/api/stories/"+story.ID, nil, reader).Code != 404 {
		t.Fatal("Non-owner deletion must fail")
	}
	if getProgress(reader).ChapterID != first.ID {
		t.Fatal("Rejected deletion destroyed reading progress")
	}
	if request("DELETE", "/api/chapters/"+first.ID, nil, author).Code != 204 {
		t.Fatal("Chapter deletion failed")
	}
	if getProgress(reader).ChapterID != "" {
		t.Fatal("Deleted chapter must clear stale progress")
	}
	if len(list()) != 0 {
		t.Fatal("Story should disappear when only one readable chapter remains")
	}
	empty := addChapter(story, 3, "<script>alert('x')</script><p> </p>")
	if request("PUT", path, map[string]string{"chapterId": empty.ID}, reader).Code != 404 {
		t.Fatal("Unreadable chapter cannot become progress")
	}
	if request("PUT", path, map[string]string{"chapterId": second.ID}, reader).Code != 204 {
		t.Fatal("Saving remaining chapter failed")
	}
	if request("DELETE", "/api/stories/"+story.ID, nil, author).Code != 204 {
		t.Fatal("Story deletion failed")
	}
	if request("GET", path, nil, reader).Code != 404 {
		t.Fatal("Deleted story must not return history")
	}
}
