package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"
	"ui-design-backend/internal/domain"
)

func TestReadingHistoryDatesMetadataAndPersistence(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "reading.db")
	db, err := NewDB("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { db.Close() }()
	now := time.Now().UTC()
	for _, id := range []string{"author", "reader", "other"} {
		if _, err := db.Exec(`INSERT INTO users(id, display_name, email, password_hash, created_at, updated_at) VALUES ($1, $1, $1, 'hash', $2, $2)`, id, now); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewStoryRepo(db)
	for _, id := range []string{"first", "single"} {
		story := &domain.Story{ID: id, UserID: "author", Title: id, Status: "published", CreatedAt: now, UpdatedAt: now}
		if err := repo.CreateStory(ctx, story); err != nil {
			t.Fatal(err)
		}
		chapter := &domain.Chapter{ID: id + "-1", StoryID: id, ChapterNumber: 1, Title: "Chương " + id, Content: "Nội dung", CreatedAt: now, UpdatedAt: now}
		if err := repo.CreateChapter(ctx, chapter); err != nil {
			t.Fatal(err)
		}
		if err := repo.SaveReadingProgress(ctx, "reader", id, chapter.ID); err != nil {
			t.Fatal(err)
		}
	}
	chapter := &domain.Chapter{ID: "first-2", StoryID: "first", ChapterNumber: 2, Title: "Chương hai", Content: "Nội dung hai", CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateChapter(ctx, chapter); err != nil {
		t.Fatal(err)
	}
	older := now.Add(-48 * time.Hour)
	recent := now.Add(-24 * time.Hour)
	for id, date := range map[string]time.Time{"first": older, "single": recent} {
		if _, err := db.Exec(`UPDATE reading_progress SET updated_at = $1 WHERE user_id = 'reader' AND story_id = $2`, date, id); err != nil {
			t.Fatal(err)
		}
	}
	// Dates and chapter metadata must survive a new database connection.
	db.Close()
	db, err = NewDB("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	repo = NewStoryRepo(db)
	reading, err := repo.ListReadingStories(ctx, "reader")
	if err != nil {
		t.Fatal(err)
	}
	if len(reading) != 2 || reading[0].ID != "single" || reading[1].ID != "first" {
		t.Fatalf("Must order by reading time and include one-chapter stories: %+v", reading)
	}
	progress := reading[0].ReadingProgress
	if progress == nil || progress.ChapterID != "single-1" || progress.ChapterNumber != 1 || progress.ChapterTitle != "Chương single" || progress.UpdatedAt == nil || !progress.UpdatedAt.Equal(recent) {
		t.Fatalf("Incorrect persisted reading metadata: %+v", progress)
	}
	public, err := repo.ListStories(ctx, "")
	if err != nil || len(public) != 1 || public[0].ID != "first" || public[0].ReadingProgress != nil {
		t.Fatalf("Public library must retain its threshold and omit personal history: %+v, %v", public, err)
	}
	other, err := repo.ListReadingStories(ctx, "other")
	if err != nil || len(other) != 0 {
		t.Fatalf("History leaked into another account: %+v, %v", other, err)
	}
	// Reopening the same chapter records a new visit, using server time.
	before := time.Now().UTC()
	if err := repo.SaveReadingProgress(ctx, "reader", "first", "first-1"); err != nil {
		t.Fatal(err)
	}
	reading, err = repo.ListReadingStories(ctx, "reader")
	if err != nil || reading[0].ID != "first" || reading[0].ReadingProgress.UpdatedAt.Before(before) {
		t.Fatalf("Rereading must move the story to the top: %+v, %v", reading, err)
	}
	// A draft can appear only in its owner's personal history.
	draft := &domain.Story{ID: "draft", UserID: "author", Title: "Bản nháp", Status: "draft", CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateStory(ctx, draft); err != nil {
		t.Fatal(err)
	}
	draftChapter := &domain.Chapter{ID: "draft-1", StoryID: "draft", ChapterNumber: 1, Title: "Bản nháp", Content: "Nội dung", CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateChapter(ctx, draftChapter); err != nil {
		t.Fatal(err)
	}
	for _, user := range []string{"reader", "author"} {
		if err := repo.SaveReadingProgress(ctx, user, "draft", "draft-1"); err != nil {
			t.Fatal(err)
		}
		items, err := repo.ListReadingStories(ctx, user)
		if err != nil {
			t.Fatal(err)
		}
		for _, item := range items {
			if item.ID == "draft" && user != "author" {
				t.Fatal("Another author's draft appeared in reading history")
			}
		}
		if user == "author" && (len(items) != 1 || items[0].ID != "draft") {
			t.Fatal("Author must retain their own draft reading history")
		}
	}
}
