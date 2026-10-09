package repository

import "ui-design-backend/internal/domain"

// Versioned and transactional so existing databases retain their stories.
func (db *DB) migrateReading() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(64) PRIMARY KEY)`); err != nil {
		return err
	}
	var count int
	if err = tx.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = 'reader-v1'`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return tx.Commit()
	}
	if _, err = tx.Exec(`ALTER TABLE chapters ADD COLUMN has_content INT NOT NULL DEFAULT 0`); err != nil {
		return err
	}
	rows, err := tx.Query(`SELECT id, COALESCE(content, '') FROM chapters`)
	if err != nil {
		return err
	}
	var readable []string
	for rows.Next() {
		var id, content string
		if err := rows.Scan(&id, &content); err != nil {
			rows.Close()
			return err
		}
		if domain.HasChapterContent(content) {
			readable = append(readable, id)
		}
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	for _, id := range readable {
		if _, err = tx.Exec(`UPDATE chapters SET has_content = 1 WHERE id = $1`, id); err != nil {
			return err
		}
	}
	if _, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_chapters_readable ON chapters(story_id, has_content);
	CREATE TABLE IF NOT EXISTS reading_progress (
		user_id VARCHAR(64) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		story_id VARCHAR(64) NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
		chapter_id VARCHAR(64) NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
		updated_at TIMESTAMP NOT NULL,
		PRIMARY KEY(user_id, story_id)
	);
	CREATE INDEX IF NOT EXISTS idx_reading_progress_chapter ON reading_progress(chapter_id);
	INSERT INTO schema_migrations(version) VALUES ('reader-v1');`); err != nil {
		return err
	}
	return tx.Commit()
}
