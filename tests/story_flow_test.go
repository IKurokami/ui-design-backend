package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"ui-design-backend/internal/auth"
	"ui-design-backend/internal/config"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/mailer"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
	"ui-design-backend/internal/service"
	transporthttp "ui-design-backend/internal/transport/http"
)

func setupStoryTestServer(t *testing.T) (http.Handler, string, func()) {
	t.Helper()

	dbFile := "./test_story_" + uuid.New().String() + ".db"
	db, err := repository.NewDB("sqlite", dbFile)
	if err != nil {
		t.Fatalf("Failed to init sqlite: %v", err)
	}

	cfg := &config.Config{
		Env:            "test",
		Port:           "8080",
		ClientOrigin:   "http://localhost:5173",
		JWTSecret:      "test-secret-at-least-32-bytes-long-super-secure!",
		AccessTokenTTL: 15 * time.Minute,
		UploadDir:      "./test_uploads_" + uuid.New().String(),
	}

	userRepo := repository.NewUserRepo(db)
	sessionRepo := repository.NewSessionRepo(db)
	resetRepo := repository.NewPasswordResetRepo(db)
	profileRepo := repository.NewProfileRepo(db)
	storyRepo := repository.NewStoryRepo(db)

	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL)
	mailService := mailer.NewMailer(cfg)
	authService := service.NewAuthService(cfg, userRepo, sessionRepo, resetRepo, tokenService, mailService)
	profileService := service.NewProfileService(userRepo, profileRepo)
	storyService := service.NewStoryService(storyRepo)

	authMW := middleware.NewAuthMiddleware(cfg, tokenService, authService)
	authHandler := transporthttp.NewAuthHandler(authService, authMW, transporthttp.NewAdminHandler(service.NewManagementService(userRepo, storyService)))
	profileHandler := transporthttp.NewProfileHandler(profileService, cfg.UploadDir)
	storyHandler := transporthttp.NewStoryHandler(storyService)

	router := transporthttp.NewRouter(cfg, authHandler, profileHandler, storyHandler, authMW)

	cleanup := func() {
		_ = db.Close()
		_ = os.Remove(dbFile)
		_ = os.RemoveAll(cfg.UploadDir)
	}

	return router, cfg.UploadDir, cleanup
}

func TestStoryAndChapterFlow(t *testing.T) {
	router, _, cleanup := setupStoryTestServer(t)
	defer cleanup()

	// 1. Register a user
	regPayload := map[string]any{
		"displayName":     "Tác giả Mộc",
		"email":           "author@example.com",
		"password":        "StrongPass123!",
		"confirmPassword": "StrongPass123!",
		"acceptedTerms":   true,
	}
	regBody, _ := json.Marshal(regPayload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Register failed: %s", w.Body.String())
	}

	var authCookie *http.Cookie
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "access_token" {
			authCookie = cookie
			break
		}
	}
	if authCookie == nil {
		t.Fatal("Missing access_token cookie")
	}

	// 2. Create a story with an initial chapter and illustrations, CoverURL left empty
	// The backend should automatically set CoverURL to the first illustration!
	createStoryReq := domain.CreateStoryDTO{
		Title:       "Đại Đạo Triều Thiên",
		Description: "Một câu chuyện tiên hiệp tu chân kỳ ảo.",
		Category:    "Tiên Hiệp",
		Tags:        "Tu Chân, Huyền Huyễn",
		Status:      "published",
		InitialChapter: &domain.CreateChapterDTO{
			ChapterNumber: 1,
			Title:         "Chương 1: Thiếu Niên Bên Suối",
			Content:       "<p>Trời xanh mây trắng...</p><p><img src=\"/uploads/cover-scene-1.jpg\" /></p>",
			Illustrations: `["/uploads/cover-scene-1.jpg", "/uploads/cover-scene-2.jpg"]`,
		},
	}
	body, _ := json.Marshal(createStoryReq)
	req = httptest.NewRequest("POST", "/api/stories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("CreateStory failed: %s", w.Body.String())
	}

	var createdStory domain.Story
	if err := json.Unmarshal(w.Body.Bytes(), &createdStory); err != nil {
		t.Fatalf("Failed to parse created story: %v", err)
	}

	if createdStory.ID == "" {
		t.Fatal("Created story ID is empty")
	}
	if createdStory.Title != "Đại Đạo Triều Thiên" {
		t.Fatalf("Expected title 'Đại Đạo Triều Thiên', got '%s'", createdStory.Title)
	}
	// Check auto-cover behavior:
	if createdStory.CoverURL != "/uploads/cover-scene-1.jpg" {
		t.Fatalf("Expected auto cover '/uploads/cover-scene-1.jpg', got '%s'", createdStory.CoverURL)
	}
	if len(createdStory.Chapters) != 1 {
		t.Fatalf("Expected 1 chapter, got %d", len(createdStory.Chapters))
	}
	ch1 := createdStory.Chapters[0]
	if ch1.Title != "Chương 1: Thiếu Niên Bên Suối" {
		t.Fatalf("Expected chapter title 'Chương 1: Thiếu Niên Bên Suối', got '%s'", ch1.Title)
	}

	// 3. Get story by ID
	req = httptest.NewRequest("GET", "/api/stories/"+createdStory.ID, nil)
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GetStory failed: %s", w.Body.String())
	}
	var fetchedStory domain.Story
	if err := json.Unmarshal(w.Body.Bytes(), &fetchedStory); err != nil {
		t.Fatalf("Failed to parse fetched story: %v", err)
	}
	if fetchedStory.AuthorName != "Tác giả Mộc" {
		t.Fatalf("Expected AuthorName 'Tác giả Mộc', got '%s'", fetchedStory.AuthorName)
	}

	// 4. List stories
	req = httptest.NewRequest("GET", "/api/stories", nil)
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("ListStories failed: %s", w.Body.String())
	}
	var stories []*domain.Story
	if err := json.Unmarshal(w.Body.Bytes(), &stories); err != nil {
		t.Fatalf("Failed to parse stories list: %v", err)
	}
	if len(stories) != 0 {
		t.Fatal("A one-chapter story must not appear in the public library")
	}

	// 5. Add chapter 2 to story
	ch2Req := domain.CreateChapterDTO{
		ChapterNumber: 2,
		Title:         "Chương 2: Kiếm Quang Khởi",
		Content:       "<p>Một đạo thanh quang chiếu rọi...</p>",
		Illustrations: `["/uploads/sword-light.jpg"]`,
	}
	body, _ = json.Marshal(ch2Req)
	req = httptest.NewRequest("POST", "/api/stories/"+createdStory.ID+"/chapters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("CreateChapter failed: %s", w.Body.String())
	}
	var createdCh2 domain.Chapter
	if err := json.Unmarshal(w.Body.Bytes(), &createdCh2); err != nil {
		t.Fatalf("Failed to parse chapter 2: %v", err)
	}
	if createdCh2.ChapterNumber != 2 {
		t.Fatalf("Expected chapter number 2, got %d", createdCh2.ChapterNumber)
	}

	// 6. Update Chapter 2 content
	newContent := "<p>Nội dung chương 2 đã được cập nhật.</p>"
	upChReq := domain.UpdateChapterDTO{
		Content: &newContent,
	}
	body, _ = json.Marshal(upChReq)
	req = httptest.NewRequest("PUT", "/api/chapters/"+createdCh2.ID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("UpdateChapter failed: %s", w.Body.String())
	}

	// 7. Update Story cover explicitly (e.g. user selected photo 2 as cover)
	newCover := "/uploads/cover-scene-2.jpg"
	upStoryReq := domain.UpdateStoryDTO{
		CoverURL: &newCover,
	}
	body, _ = json.Marshal(upStoryReq)
	req = httptest.NewRequest("PUT", "/api/stories/"+createdStory.ID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("UpdateStory failed: %s", w.Body.String())
	}
	var updatedStory domain.Story
	if err := json.Unmarshal(w.Body.Bytes(), &updatedStory); err != nil {
		t.Fatalf("Failed to parse updated story: %v", err)
	}
	if updatedStory.CoverURL != newCover {
		t.Fatalf("Expected cover '%s', got '%s'", newCover, updatedStory.CoverURL)
	}

	// 8. Tủ sách cá nhân (Bookshelf): GET /api/stories/mine
	req = httptest.NewRequest("GET", "/api/stories/mine", nil)
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/stories/mine failed: %d - %s", w.Code, w.Body.String())
	}
	var myStories []domain.Story
	if err := json.Unmarshal(w.Body.Bytes(), &myStories); err != nil {
		t.Fatalf("Failed to parse bookshelf stories: %v", err)
	}
	if len(myStories) != 1 {
		t.Fatalf("Expected 1 story in my bookshelf, got %d", len(myStories))
	}
	if myStories[0].ChapterCount != 2 {
		t.Fatalf("Expected ChapterCount 2 in bookshelf, got %d", myStories[0].ChapterCount)
	}
	expectedPreviews := []string{newCover, "/uploads/cover-scene-1.jpg", "/uploads/sword-light.jpg"}
	if !reflect.DeepEqual(myStories[0].PreviewImages, expectedPreviews) {
		t.Fatalf("Incorrect bookshelf previews: %v", myStories[0].PreviewImages)
	}
	// The public list returns distinct images without sending full chapter contents.
	req = httptest.NewRequest("GET", "/api/stories", nil)
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("ListStories failed: %s", w.Body.String())
	}
	if err := json.Unmarshal(w.Body.Bytes(), &stories); err != nil {
		t.Fatal(err)
	}
	if len(stories) != 1 || !reflect.DeepEqual(stories[0].PreviewImages, expectedPreviews) {
		t.Fatalf("Incorrect public previews: %s", w.Body.String())
	}
	if len(stories[0].Chapters) != 0 {
		t.Fatal("List response must not include full chapters")
	}

	// 9. Unauthenticated access to /api/stories/mine must fail with 401
	unauthReq := httptest.NewRequest("GET", "/api/stories/mine", nil)
	unauthW := httptest.NewRecorder()
	router.ServeHTTP(unauthW, unauthReq)
	if unauthW.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for unauthenticated /api/stories/mine, got %d", unauthW.Code)
	}

	// 10. Delete Story from Bookshelf: DELETE /api/stories/{id}
	delReq := httptest.NewRequest("DELETE", "/api/stories/"+createdStory.ID, nil)
	delReq.AddCookie(authCookie)
	delW := httptest.NewRecorder()
	router.ServeHTTP(delW, delReq)

	if delW.Code != http.StatusOK && delW.Code != http.StatusNoContent {
		t.Fatalf("DELETE /api/stories failed: %d - %s", delW.Code, delW.Body.String())
	}

	// Verify bookshelf is now empty
	req = httptest.NewRequest("GET", "/api/stories/mine", nil)
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if err := json.Unmarshal(w.Body.Bytes(), &myStories); err != nil {
		t.Fatalf("Failed to parse bookshelf stories after delete: %v", err)
	}
	if len(myStories) != 0 {
		t.Fatalf("Expected 0 stories in bookshelf after delete, got %d", len(myStories))
	}
}
