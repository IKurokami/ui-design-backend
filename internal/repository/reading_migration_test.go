package repository

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestReadingMigrationBackfillsExistingChaptersAndIsRepeatable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.db")
	legacy, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = legacy.Exec(`CREATE TABLE chapters (
 id VARCHAR(64) PRIMARY KEY, story_id VARCHAR(64), chapter_number INT, title TEXT, content TEXT,
 illustrations TEXT, created_at TIMESTAMP, updated_at TIMESTAMP);
 INSERT INTO chapters(id, story_id, content) VALUES ('readable', 'story', '<p>Nội dung</p>'), ('empty', 'story', '<p><br>&nbsp;</p>');`)
	legacy.Close()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		db, err := NewDB("sqlite", path)
		if err != nil {
			t.Fatal(err)
		}
		var readable, empty int
		if err := db.QueryRow(`SELECT has_content FROM chapters WHERE id = 'readable'`).Scan(&readable); err != nil {
			t.Fatal(err)
		}
		if err := db.QueryRow(`SELECT has_content FROM chapters WHERE id = 'empty'`).Scan(&empty); err != nil {
			t.Fatal(err)
		}
		if readable != 1 || empty != 0 {
			t.Fatal("Existing content was not backfilled correctly")
		}
		db.Close()
	}
}
