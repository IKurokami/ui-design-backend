package repository

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ui-design-backend/internal/domain"
)

func TestRetiredPresentationMigrationPreservesChaptersAndHistory(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "upgrade.db")
	db, err := NewDB("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	for _, id := range []string{"author", "reader"} {
		if _, err = db.Exec(`INSERT INTO users(id, display_name, email, password_hash, created_at, updated_at) VALUES ($1, $1, $1, 'hash', $2, $2)`, id, now); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewStoryRepo(db)
	story := &domain.Story{ID: "story", UserID: "author", Title: "Truyện", CoverURL: "/uploads/cover.jpg", Status: "published", CreatedAt: now, UpdatedAt: now}
	if err = repo.CreateStory(ctx, story); err != nil {
		t.Fatal(err)
	}
	content := `<p><em style="color: red">Định dạng chữ</em></p><div data-book-tiles="2" style="display:grid"><figure data-book-frame="true" contenteditable="false"><img src="/uploads/scene.jpg" data-book-width="40" data-book-zoom="1.8" style="transform:scale(1.8)"></figure><div data-book-tile-text="true"><p>Chữ bên cạnh ảnh</p></div></div>`
	if err = repo.CreateChapter(ctx, &domain.Chapter{ID: "chapter", StoryID: story.ID, Title: "Chương một", ChapterNumber: 1, Content: content, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	if err = repo.SaveReadingProgress(ctx, "reader", story.ID, "chapter"); err != nil {
		t.Fatal(err)
	}
	original, err := repo.GetReadingProgress(ctx, "reader", story.ID)
	if err != nil {
		t.Fatal(err)
	}
	// Reproduce a database upgraded by the retired feature, rather than just a fresh schema.
	for _, query := range []string{
		`ALTER TABLE stories ADD COLUMN book_layout TEXT NOT NULL DEFAULT '{}'`,
		`ALTER TABLE reading_progress ADD COLUMN mode VARCHAR(16) NOT NULL DEFAULT 'book'`,
		`ALTER TABLE reading_progress ADD COLUMN page_index INT NOT NULL DEFAULT 10`,
		`ALTER TABLE reading_progress ADD COLUMN block_index INT NOT NULL DEFAULT 2`,
		`ALTER TABLE reading_progress ADD COLUMN text_offset INT NOT NULL DEFAULT 50`,
		`ALTER TABLE reading_progress ADD COLUMN scroll_progress DOUBLE PRECISION NOT NULL DEFAULT 0`,
		`DELETE FROM schema_migrations WHERE version = 'reader-chapter-only-v1'`,
		`INSERT INTO schema_migrations(version) VALUES ('reader-book-v1')`,
	} {
		if _, err = db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		db, err = NewDB("sqlite", path)
		if err != nil {
			t.Fatal(err)
		}
		repo = NewStoryRepo(db)
		saved, err := repo.GetStory(ctx, story.ID)
		if err != nil {
			t.Fatal(err)
		}
		if saved.Title != story.Title || saved.CoverURL != story.CoverURL || len(saved.Chapters) != 1 {
			t.Fatal("Migration changed story data")
		}
		html := saved.Chapters[0].Content
		for _, preserved := range []string{"/uploads/scene.jpg", "Chữ bên cạnh ảnh", `<em style="color: red">Định dạng chữ</em>`} {
			if !strings.Contains(html, preserved) {
				t.Fatalf("Lost chapter content: %s", html)
			}
		}
		for _, removed := range []string{"data-book-", "contenteditable", "transform:", "display:grid", "<figure"} {
			if strings.Contains(html, removed) {
				t.Fatalf("Retired layout remains: %s", html)
			}
		}
		progress, err := repo.GetReadingProgress(ctx, "reader", story.ID)
		if err != nil || progress.ChapterID != original.ChapterID || !progress.UpdatedAt.Equal(*original.UpdatedAt) {
			t.Fatalf("Lost reading history: %+v, %v", progress, err)
		}
		history, err := repo.ListReadingStories(ctx, "reader")
		if err != nil || len(history) != 1 || history[0].ReadingProgress.ChapterID != "chapter" {
			t.Fatal("Reading list must remain intact")
		}
		tx, err := db.Begin()
		if err != nil {
			t.Fatal(err)
		}
		for _, table := range []string{"stories", "reading_progress"} {
			columns, err := schemaColumns(tx, "sqlite", table)
			if err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"book_layout", "mode", "page_index", "block_index", "text_offset", "scroll_progress"} {
				if columns[name] {
					t.Fatalf("Retired column %s.%s remains", table, name)
				}
			}
		}
		tx.Rollback()
		if err = repo.SaveReadingProgress(ctx, "reader", story.ID, "chapter"); err != nil {
			t.Fatal(err)
		}
		original, err = repo.GetReadingProgress(ctx, "reader", story.ID)
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}
