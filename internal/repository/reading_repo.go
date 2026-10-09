package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"ui-design-backend/internal/domain"
)

func contentFlag(content string) int {
	if domain.HasChapterContent(content) {
		return 1
	}
	return 0
}

func (r *StoryRepo) GetReadingProgress(ctx context.Context, userID, storyID string) (*domain.ReadingProgress, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM stories WHERE id = $1)`, storyID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrStoryNotFound
	}
	var progress domain.ReadingProgress
	var updated time.Time
	err := r.db.QueryRowContext(ctx, `SELECT p.chapter_id, c.chapter_number, c.title, p.updated_at FROM reading_progress p
	JOIN chapters c ON c.id = p.chapter_id AND c.story_id = p.story_id AND c.has_content = 1
	JOIN stories s ON s.id = p.story_id
	WHERE p.user_id = $1 AND p.story_id = $2`, userID, storyID).Scan(&progress.ChapterID, &progress.ChapterNumber, &progress.ChapterTitle, &updated)
	if errors.Is(err, sql.ErrNoRows) {
		return &progress, nil
	}
	if err != nil {
		return nil, err
	}
	progress.UpdatedAt = &updated
	return &progress, nil
}

func (r *StoryRepo) SaveReadingProgress(ctx context.Context, userID, storyID, chapterID string) error {
	// Validate and write in one statement; a deleted/wrong/empty chapter cannot become progress.
	res, err := r.db.ExecContext(ctx, `INSERT INTO reading_progress(user_id, story_id, chapter_id, updated_at)
	SELECT $1, c.story_id, c.id, $4 FROM chapters c WHERE c.id = $3 AND c.story_id = $2 AND c.has_content = 1
	ON CONFLICT(user_id, story_id) DO UPDATE SET chapter_id = excluded.chapter_id, updated_at = excluded.updated_at`,
		userID, storyID, chapterID, time.Now().UTC())
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrChapterNotFound
	}
	return nil
}
