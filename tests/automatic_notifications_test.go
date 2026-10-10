package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "ui-design-backend/internal/domain"
)

func TestAutomaticChapterNotificationsFromReadingAndFavorites(t *testing.T) {
    router, _, cleanup := setupStoryTestServer(t)
    defer cleanup()

    request := func(method, path string, payload any, cookie *http.Cookie) *httptest.ResponseRecorder {
        var body bytes.Buffer
        if payload != nil {
            if err := json.NewEncoder(&body).Encode(payload); err != nil { t.Fatal(err) }
        }
        req := httptest.NewRequest(method, path, &body)
        req.Header.Set("Content-Type", "application/json")
        if cookie != nil { req.AddCookie(cookie) }
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        return w
    }
    register := func(email string) *http.Cookie {
        w := request("POST", "/api/auth/register", map[string]any{
            "displayName": "Độc giả", "email": email, "password": "StrongPass123!",
            "confirmPassword": "StrongPass123!", "acceptedTerms": true,
        }, nil)
        if w.Code != 200 { t.Fatalf("register %s failed: %d: %s", email, w.Code, w.Body.String()) }
        for _, cookie := range w.Result().Cookies() {
            if cookie.Name == "access_token" { return cookie }
        }
        t.Fatalf("missing access token for %s", email)
        return nil
    }
    author := register("automatic-author@example.com")
    history := register("automatic-history@example.com")
    favorite := register("automatic-favorite@example.com")
    follower := register("automatic-follower@example.com")
    overlap := register("automatic-overlap@example.com")
    muted := register("automatic-muted@example.com")
    stranger := register("automatic-stranger@example.com")

    w := request("POST", "/api/stories", domain.CreateStoryDTO{Title: "Truyện quan tâm", Status: "published"}, author)
    if w.Code != 201 { t.Fatalf("create story: %d %s", w.Code, w.Body.String()) }
    var story domain.Story
    if err := json.Unmarshal(w.Body.Bytes(), &story); err != nil { t.Fatal(err) }

    newChapter := func(title, content string) domain.Chapter {
        w := request("POST", "/api/stories/"+story.ID+"/chapters",
            domain.CreateChapterDTO{Title: title, Content: content}, author)
        if w.Code != 201 { t.Fatalf("create chapter: %d %s", w.Code, w.Body.String()) }
        var chapter domain.Chapter
        if err := json.Unmarshal(w.Body.Bytes(), &chapter); err != nil { t.Fatal(err) }
        return chapter
    }
    first := newChapter("Chương 1", "<p>Nội dung 1</p>")
    newChapter("Chương 2", "<p>Nội dung 2</p>")

    progressPath := "/api/stories/" + story.ID + "/reading-progress"
    for _, user := range []*http.Cookie{history, overlap, author} {
        if got := request("PUT", progressPath, map[string]string{"chapterId": first.ID}, user).Code; got != 204 {
            t.Fatalf("save reading progress: %d", got)
        }
    }
    for _, user := range []*http.Cookie{favorite, overlap, muted} {
        if got := request("POST", "/api/favorites/"+story.ID, nil, user).Code; got != 204 {
            t.Fatalf("save story: %d", got)
        }
    }
    for _, user := range []*http.Cookie{follower, overlap} {
        if got := request("POST", "/api/follows/"+story.ID, nil, user).Code; got != 204 {
            t.Fatalf("follow story: %d", got)
        }
    }
    if got := request("POST", "/api/notifications/mute/"+story.ID, nil, nil).Code; got != 401 {
        t.Fatalf("anonymous mute: %d", got)
    }
    if got := request("POST", "/api/notifications/mute/"+story.ID, nil, muted).Code; got != 204 {
        t.Fatalf("mute: %d", got)
    }
    if got := request("GET", "/api/notifications/muted", nil, stranger).Code; got != 200 { t.Fatal(got) }
    w = request("GET", "/api/notifications/muted", nil, muted)
    var muteIDs struct { IDs []string `json:"ids"` }
    if err := json.Unmarshal(w.Body.Bytes(), &muteIDs); err != nil { t.Fatal(err) }
    if len(muteIDs.IDs) != 1 || muteIDs.IDs[0] != story.ID {
        t.Fatalf("mute did not persist: %+v", muteIDs.IDs)
    }

    third := newChapter("Chương 3", "<p>Nội dung 3</p>")
    loadInbox := func(user *http.Cookie) []domain.StoryNotification {
        w := request("GET", "/api/notifications", nil, user)
        if w.Code != 200 { t.Fatalf("inbox response: %d %s", w.Code, w.Body.String()) }
        var result struct { Items []domain.StoryNotification `json:"items"` }
        if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil { t.Fatal(err) }
        return result.Items
    }
    for name, user := range map[string]*http.Cookie{
        "reading history": history, "favorite": favorite, "follow": follower, "all signals": overlap,
    } {
        inbox := loadInbox(user)
        if len(inbox) != 1 || inbox[0].ChapterID != third.ID || inbox[0].Read {
            t.Errorf("%s: expected exactly one unread notification, got %+v", name, inbox)
        }
    }
    for name, user := range map[string]*http.Cookie{
        "muted": muted, "unrelated": stranger, "author": author,
    } {
        if inbox := loadInbox(user); len(inbox) != 0 {
            t.Errorf("%s received unrelated notifications: %+v", name, inbox)
        }
    }

    // Editing an already-readable chapter never creates a second alert.
    if got := request("PUT", "/api/chapters/"+third.ID, domain.UpdateChapterDTO{
        Content: ptrString("<p>Đã sửa nội dung</p>"),
    }, author).Code; got != 200 { t.Fatalf("edit chapter: %d", got) }
    if len(loadInbox(overlap)) != 1 { t.Fatal("duplicate alert after editing chapter content") }

    // A chapter initially without readable content should notify on its first publication.
    empty := newChapter("Chương 4 (nháp trống)", "<p> </p>")
    if len(loadInbox(favorite)) != 1 { t.Fatal("blank chapter triggered notification") }
    if got := request("PUT", "/api/chapters/"+empty.ID,
        domain.UpdateChapterDTO{Content: ptrString("<p>Nội dung mới công khai</p>")}, author).Code; got != 200 {
        t.Fatalf("publish blank chapter: %d", got)
    }
    if len(loadInbox(favorite)) != 2 { t.Fatal("publishing readable content did not create alert") }
    if len(loadInbox(muted)) != 0 { t.Fatal("muted user received chapter notification") }

    // Unmute without removing the saved story. Future chapters notify again.
    if got := request("DELETE", "/api/notifications/mute/"+story.ID, nil, muted).Code; got != 204 { t.Fatal(got) }
    next := newChapter("Chương 5", "<p>Chương 5</p>")
    after := loadInbox(muted)
    if len(after) != 1 || after[0].ChapterID != next.ID { t.Fatalf("unmute didn't restore notification: %+v", after) }
}

func ptrString(s string) *string { return &s }
