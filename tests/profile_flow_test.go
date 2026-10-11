package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ui-design-backend/internal/domain"
)

func TestProfileAndUploadFlow(t *testing.T) {
	router, _, cleanup := setupTestApp(t)
	defer cleanup()
	defer os.RemoveAll("./test_uploads")
	defer os.RemoveAll("./uploads")

	// 1. Register a test user
	registerBody := `{"displayName":"Designer Pro","email":"designer@example.com","password":"Password123","confirmPassword":"Password123","acceptedTerms":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(registerBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Register failed: status %d, body: %s", w.Code, w.Body.String())
	}

	var sessionResp domain.AuthSessionResponse
	_ = json.Unmarshal(w.Body.Bytes(), &sessionResp)
	userID := sessionResp.User.ID

	var accessTokenCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "access_token" {
			accessTokenCookie = c
			break
		}
	}
	if accessTokenCookie == nil {
		t.Fatalf("Missing access_token cookie")
	}

	// 2. Upload an avatar image file via POST /api/upload
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	part, err := writer.CreateFormFile("file", "avatar_sample.png")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	sampleImageBytes := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDRsampleimagecontent")
	_, _ = part.Write(sampleImageBytes)
	writer.Close()

	uploadReq := httptest.NewRequest(http.MethodPost, "/api/upload", bodyBuf)
	uploadReq.Header.Set("Content-Type", writer.FormDataContentType())
	uploadReq.AddCookie(accessTokenCookie)
	uploadRec := httptest.NewRecorder()
	router.ServeHTTP(uploadRec, uploadReq)

	if uploadRec.Code != http.StatusOK {
		t.Fatalf("Upload failed with status %d: %s", uploadRec.Code, uploadRec.Body.String())
	}

	var uploadResp struct {
		URL      string `json:"url"`
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
	}
	if err := json.Unmarshal(uploadRec.Body.Bytes(), &uploadResp); err != nil {
		t.Fatalf("Failed to parse upload response: %v", err)
	}

	if uploadResp.URL == "" || uploadResp.Filename == "" {
		t.Fatalf("Expected non-empty URL and Filename, got %+v", uploadResp)
	}

	// 3. Get /api/profile/me (should auto-create default profile with #3B82F6)
	getMeReq := httptest.NewRequest(http.MethodGet, "/api/profile/me", nil)
	getMeReq.AddCookie(accessTokenCookie)
	getMeRec := httptest.NewRecorder()
	router.ServeHTTP(getMeRec, getMeReq)

	if getMeRec.Code != http.StatusOK {
		t.Fatalf("GetMe profile failed: %d, body: %s", getMeRec.Code, getMeRec.Body.String())
	}

	var myProfile domain.ProfileResponse
	if err := json.Unmarshal(getMeRec.Body.Bytes(), &myProfile); err != nil {
		t.Fatalf("Failed to decode profile response: %v", err)
	}

	if myProfile.DisplayName != "Designer Pro" || myProfile.ThemeColor != "#3B82F6" {
		t.Errorf("Unexpected default profile data: %+v", myProfile)
	}

	// 4. Update profile via PUT /api/profile/me (set avatar, banner, bio, headline, themeColor, etc.)
	avatarURL := uploadResp.URL
	bannerURL := "https://example.com/banner.jpg"
	bio := "Chuyên thiết kế UI/UX và hệ thống Design System."
	headline := "Lead UI/UX Designer"
	themeColor := "#10B981"
	websiteURL := "https://myportfolio.design"
	location := "TP. Hồ Chí Minh, VN"
	socialLinks := `{"github":"https://github.com/myuser","behance":"https://behance.net/myuser"}`

	updatePayload := map[string]any{
		"displayName": "Designer Pro Max",
		"avatarUrl":   avatarURL,
		"bannerUrl":   bannerURL,
		"bio":         bio,
		"headline":    headline,
		"themeColor":  themeColor,
		"websiteUrl":  websiteURL,
		"location":    location,
		"socialLinks": socialLinks,
	}
	updateJSON, _ := json.Marshal(updatePayload)

	updateReq := httptest.NewRequest(http.MethodPut, "/api/profile/me", bytes.NewBuffer(updateJSON))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.AddCookie(accessTokenCookie)
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("Update profile failed: %d, body: %s", updateRec.Code, updateRec.Body.String())
	}

	var updatedProfile domain.ProfileResponse
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updatedProfile); err != nil {
		t.Fatalf("Failed to decode updated profile: %v", err)
	}

	if updatedProfile.DisplayName != "Designer Pro Max" {
		t.Errorf("Expected updated display name 'Designer Pro Max', got %s", updatedProfile.DisplayName)
	}
	if updatedProfile.AvatarURL == nil || *updatedProfile.AvatarURL != avatarURL {
		t.Errorf("Expected avatarUrl %s, got %v", avatarURL, updatedProfile.AvatarURL)
	}
	if updatedProfile.Bio == nil || *updatedProfile.Bio != bio {
		t.Errorf("Expected bio %s, got %v", bio, updatedProfile.Bio)
	}
	if updatedProfile.ThemeColor != themeColor {
		t.Errorf("Expected themeColor %s, got %s", themeColor, updatedProfile.ThemeColor)
	}

	// 5. Profiles require login, including reads of another member's profile.
	publicReq := httptest.NewRequest(http.MethodGet, "/api/profile/"+userID, nil)
	publicRec := httptest.NewRecorder()
	router.ServeHTTP(publicRec, publicReq)
	if publicRec.Code != http.StatusUnauthorized {
		t.Fatalf("Anonymous profile request must be 401, got %d", publicRec.Code)
	}
	publicReq = httptest.NewRequest(http.MethodGet, "/api/profile/"+userID, nil)
	publicReq.AddCookie(accessTokenCookie)
	publicRec = httptest.NewRecorder()
	router.ServeHTTP(publicRec, publicReq)
	if publicRec.Code != http.StatusOK {
		t.Fatalf("Authenticated profile request failed: %d, body: %s", publicRec.Code, publicRec.Body.String())
	}

	var publicProfile domain.ProfileResponse
	if err := json.Unmarshal(publicRec.Body.Bytes(), &publicProfile); err != nil {
		t.Fatalf("Failed to decode public profile: %v", err)
	}

	if publicProfile.UserID != userID || publicProfile.Email != "" {
		t.Errorf("Public profile should not expose email, got email: %s", publicProfile.Email)
	}
	if publicProfile.Headline == nil || *publicProfile.Headline != headline {
		t.Errorf("Expected headline %s, got %v", headline, publicProfile.Headline)
	}

	// 6. Upload URLs must not disclose assets to logged-out visitors.
	staticReq := httptest.NewRequest(http.MethodGet, uploadResp.URL, nil)
	staticRec := httptest.NewRecorder()
	router.ServeHTTP(staticRec, staticReq)
	if staticRec.Code != http.StatusUnauthorized {
		t.Fatalf("Anonymous file download must be 401, got %d", staticRec.Code)
	}
	staticReq = httptest.NewRequest(http.MethodGet, uploadResp.URL, nil)
	staticReq.AddCookie(accessTokenCookie)
	staticRec = httptest.NewRecorder()
	router.ServeHTTP(staticRec, staticReq)
	if staticRec.Code != http.StatusOK {
		t.Fatalf("Authenticated file retrieval failed: status %d", staticRec.Code)
	}
	if staticRec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("Uploaded image responses must disable MIME sniffing")
	}
	body, _ := io.ReadAll(staticRec.Body)
	if len(body) != len(sampleImageBytes) {
		t.Errorf("Static file content length mismatch: got %d, expected %d", len(body), len(sampleImageBytes))
	}
}

func TestRejectActiveOrMisleadingImageUploads(t *testing.T) {
	router, _, cleanup := setupTestApp(t)
	defer cleanup()
	registration := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(
		`{"displayName":"Image checker","email":"safe-image-check@example.com","password":"Password123","confirmPassword":"Password123","acceptedTerms":true}`))
	registration.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, registration)
	if rw.Code != http.StatusOK {
		t.Fatalf("register: %d %s", rw.Code, rw.Body.String())
	}
	var cookie *http.Cookie
	for _, c := range rw.Result().Cookies() {
		if c.Name == "access_token" {
			cookie = c
		}
	}
	if cookie == nil {
		t.Fatal("missing auth cookie")
	}
	for _, testCase := range []struct {
		name string
		filename string
		content string
	}{
		{"svg with active script", "attack.svg", `<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`},
		{"HTML renamed to PNG", "attack.png", `<script>alert(1)</script>`},
		{"PNG renamed to JPEG", "attack.jpg", "\x89PNG\r\n\x1a\ncontent"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", testCase.filename)
			if err != nil { t.Fatal(err) }
			if _, err := part.Write([]byte(testCase.content)); err != nil { t.Fatal(err) }
			if err := writer.Close(); err != nil { t.Fatal(err) }
			req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.AddCookie(cookie)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("unsafe upload %s accepted: %d %s", testCase.filename, w.Code, w.Body.String())
			}
		})
	}
}
