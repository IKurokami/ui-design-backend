package repository

import (
	"context"
	"time"

	"ui-design-backend/internal/domain"
)

// Favorites are scoped to a user. Only stories visible in the library
// (published/completed with at least two readable chapters) may be saved.
const visibleFavoriteSQL = "s.status IN ('published', 'completed') AND (SELECT COUNT(*) FROM chapters c WHERE c.story_id = s.id AND c.has_content = 1) > 1"

func (r *StoryRepo) FavoriteIDs(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT f.story_id FROM favorite_stories f
		JOIN stories s ON s.id = f.story_id WHERE f.user_id = $1 AND `+visibleFavoriteSQL+`
		ORDER BY f.created_at DESC, f.story_id ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *StoryRepo) ListFavorites(ctx context.Context, userID string, page, size int) ([]*domain.Story, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 24
	}
	if size > 60 {
		size = 60
	}
	rows, err := r.db.QueryContext(ctx, `SELECT f.story_id FROM favorite_stories f
		JOIN stories s ON s.id = f.story_id WHERE f.user_id = $1 AND `+visibleFavoriteSQL+`
		ORDER BY f.created_at DESC, f.story_id ASC LIMIT $2 OFFSET $3`, userID, size, (page-1)*size)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	items := make([]*domain.Story, 0, len(ids))
	for _, id := range ids {
		story, err := r.GetStory(ctx, id)
		if err != nil {
			return nil, err
		}
		items = append(items, story)
	}
	return items, nil
}

func (r *StoryRepo) SaveFavorite(ctx context.Context, userID, storyID string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `INSERT INTO favorite_stories(user_id, story_id, created_at)
		SELECT $1, s.id, $3 FROM stories s WHERE s.id = $2 AND `+visibleFavoriteSQL+`
		ON CONFLICT(user_id, story_id) DO NOTHING`, userID, storyID, time.Now().UTC())
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected > 0 {
		return true, nil
	}
	// Idempotent duplicate save, but unknown/private stories must return 404.
	var exists int
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM favorite_stories f JOIN stories s ON f.story_id = s.id
		WHERE f.user_id = $1 AND f.story_id = $2 AND `+visibleFavoriteSQL, userID, storyID).Scan(&exists)
	return exists > 0, err
}

func (r *StoryRepo) RemoveFavorite(ctx context.Context, userID, storyID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM favorite_stories WHERE user_id = $1 AND story_id = $2", userID, storyID)
	return err
}
