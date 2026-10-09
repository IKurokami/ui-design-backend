package repository

import (
	"context"
	"database/sql"
	"errors"

	"ui-design-backend/internal/domain"
)

var (
	ErrStoryNotFound   = errors.New("story not found")
	ErrChapterNotFound = errors.New("chapter not found")
)

type StoryRepo struct {
	db *DB
}

func NewStoryRepo(db *DB) *StoryRepo {
	return &StoryRepo{db: db}
}

func (r *StoryRepo) CreateStory(ctx context.Context, s *domain.Story) error {
	query := `INSERT INTO stories (id, user_id, title, description, cover_url, category, tags, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	if r.db.Driver == "sqlite" {
		query = `INSERT INTO stories (id, user_id, title, description, cover_url, category, tags, status, created_at, updated_at)
		         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := r.db.ExecContext(ctx, query,
		s.ID,
		s.UserID,
		s.Title,
		s.Description,
		s.CoverURL,
		s.Category,
		s.Tags,
		s.Status,
		s.CreatedAt,
		s.UpdatedAt,
	)
	return err
}

func (r *StoryRepo) GetStory(ctx context.Context, id string) (*domain.Story, error) {
	query := `SELECT s.id, s.user_id, COALESCE(u.display_name, ''), s.title, s.description, s.cover_url, s.category, s.tags, s.status, s.created_at, s.updated_at
	          FROM stories s
	          LEFT JOIN users u ON s.user_id = u.id
	          WHERE s.id = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT s.id, s.user_id, COALESCE(u.display_name, ''), s.title, s.description, s.cover_url, s.category, s.tags, s.status, s.created_at, s.updated_at
		         FROM stories s
		         LEFT JOIN users u ON s.user_id = u.id
		         WHERE s.id = ?`
	}

	var s domain.Story
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.AuthorName,
		&s.Title,
		&s.Description,
		&s.CoverURL,
		&s.Category,
		&s.Tags,
		&s.Status,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStoryNotFound
		}
		return nil, err
	}

	chapters, err := r.ListChaptersByStory(ctx, s.ID)
	if err != nil {
		return nil, err
	}
	if chapters != nil {
		for _, ch := range chapters {
			s.Chapters = append(s.Chapters, *ch)
		}
	}
	s.ChapterCount = len(s.Chapters)
	s.RefreshPreviewImages()

	return &s, nil
}

func (r *StoryRepo) ListStories(ctx context.Context, userID string) ([]*domain.Story, error) {
	return r.listStories(ctx, userID, "")
}
func (r *StoryRepo) ListAll(ctx context.Context) ([]*domain.Story, error) {
	return r.listStories(ctx, "__all_management__", "")
}
func (r *StoryRepo) ListAllPage(ctx context.Context, page, pageSize int, query string) ([]*domain.Story, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	like := "%" + query + "%"
	where := " WHERE (s.title LIKE $1 OR COALESCE(u.display_name,'') LIKE $1 OR s.category LIKE $1)"
	args := []any{like}
	if r.db.Driver == "sqlite" {
		where = " WHERE (s.title LIKE ? OR COALESCE(u.display_name,'') LIKE ? OR s.category LIKE ?)"
		args = []any{like, like, like}
	}
	countQ := "SELECT COUNT(*) FROM stories s LEFT JOIN users u ON s.user_id=u.id" + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	base := `SELECT s.id,s.user_id,COALESCE(u.display_name,''),s.title,s.description,s.cover_url,s.category,s.tags,s.status,s.created_at,s.updated_at,(SELECT COUNT(*) FROM chapters c WHERE c.story_id=s.id) FROM stories s LEFT JOIN users u ON s.user_id=u.id` + where + " ORDER BY s.updated_at DESC,s.created_at DESC,s.id ASC LIMIT " + placeholder(r.db.Driver, len(args)+1) + " OFFSET " + placeholder(r.db.Driver, len(args)+2)
	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]*domain.Story, 0, pageSize)
	for rows.Next() {
		var s domain.Story
		if err := rows.Scan(&s.ID, &s.UserID, &s.AuthorName, &s.Title, &s.Description, &s.CoverURL, &s.Category, &s.Tags, &s.Status, &s.CreatedAt, &s.UpdatedAt, &s.ChapterCount); err != nil {
			return nil, 0, err
		}
		items = append(items, &s)
	}
	return items, total, rows.Err()
}
func (r *StoryRepo) ListReadingStories(ctx context.Context, userID string) ([]*domain.Story, error) {
	return r.listStories(ctx, "", userID)
}
func (r *StoryRepo) ListPage(ctx context.Context, userID, readerID string, page, pageSize int) ([]*domain.Story, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 24
	}
	if pageSize > 60 {
		pageSize = 60
	}
	return r.listStoriesPage(ctx, userID, readerID, page, pageSize)
}
func (r *StoryRepo) listStories(ctx context.Context, userID, readerID string) ([]*domain.Story, error) {
	return r.listStoriesPage(ctx, userID, readerID, 0, 0)
}
func (r *StoryRepo) listStoriesPage(ctx context.Context, userID, readerID string, page, pageSize int) ([]*domain.Story, error) {
	query := `SELECT s.id, s.user_id, COALESCE(u.display_name, ''), s.title, s.description, s.cover_url, s.category, s.tags, s.status, s.created_at, s.updated_at,
 (SELECT COUNT(*) FROM chapters c WHERE c.story_id = s.id) AS chapter_count`
	if readerID != "" {
		query += `, p.chapter_id, c.chapter_number, c.title, p.updated_at`
	}
	query += ` FROM stories s LEFT JOIN users u ON s.user_id = u.id`
	var args []any
	if readerID != "" {
		// Reading history is personal and independent of the public library's chapter threshold.
		query += ` JOIN reading_progress p ON p.story_id = s.id AND p.user_id = $1
		JOIN chapters c ON c.id = p.chapter_id AND c.story_id = s.id AND c.has_content = 1
		WHERE (s.status IN ('published', 'completed') OR s.user_id = $1)
		ORDER BY p.updated_at DESC, s.id ASC`
		args = append(args, readerID)
	} else if userID == "__all_management__" {
		// Management list includes drafts and stories with fewer than two chapters.
	} else if userID != "" {
		query += " WHERE s.user_id = $1"
		args = append(args, userID)
	} else {
		query += ` WHERE s.status IN ('published', 'completed') AND (SELECT COUNT(*) FROM chapters c WHERE c.story_id = s.id AND c.has_content = 1) > 1`
	}
	if readerID == "" {
		query += " ORDER BY s.updated_at DESC, s.created_at DESC"
	}
	if page > 0 {
		query += " LIMIT " + placeholder(r.db.Driver, len(args)+1) + " OFFSET " + placeholder(r.db.Driver, len(args)+2)
		args = append(args, pageSize, (page-1)*pageSize)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []*domain.Story
	for rows.Next() {
		var s domain.Story
		fields := []any{
			&s.ID,
			&s.UserID,
			&s.AuthorName,
			&s.Title,
			&s.Description,
			&s.CoverURL,
			&s.Category,
			&s.Tags,
			&s.Status,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.ChapterCount,
		}
		if readerID != "" {
			s.ReadingProgress = &domain.ReadingProgress{}
			fields = append(fields, &s.ReadingProgress.ChapterID, &s.ReadingProgress.ChapterNumber, &s.ReadingProgress.ChapterTitle, &s.ReadingProgress.UpdatedAt)
		}
		if err := rows.Scan(fields...); err != nil {
			return nil, err
		}
		stories = append(stories, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := r.populatePreviewImages(ctx, stories); err != nil {
		return nil, err
	}
	return stories, nil
}

func (r *StoryRepo) UpdateStory(ctx context.Context, s *domain.Story) error {
	query := `UPDATE stories
	          SET title = $1, description = $2, cover_url = $3, category = $4, tags = $5, status = $6, updated_at = $7
	          WHERE id = $8 AND user_id = $9`

	res, err := r.db.ExecContext(ctx, query,
		s.Title,
		s.Description,
		s.CoverURL,
		s.Category,
		s.Tags,
		s.Status,
		s.UpdatedAt,
		s.ID,
		s.UserID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrStoryNotFound
	}
	return nil
}

func (r *StoryRepo) DeleteStory(ctx context.Context, id, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var owner string
	if err := tx.QueryRowContext(ctx, "SELECT user_id FROM stories WHERE id = $1", id).Scan(&owner); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrStoryNotFound
		}
		return err
	}
	if owner != userID {
		return ErrStoryNotFound
	}
	for _, query := range []string{
		"DELETE FROM reading_progress WHERE story_id = $1",
		"DELETE FROM chapters WHERE story_id = $1",
		"DELETE FROM stories WHERE id = $1 AND user_id = $2",
	} {
		args := []any{id}
		if query == "DELETE FROM stories WHERE id = $1 AND user_id = $2" {
			args = append(args, userID)
		}
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *StoryRepo) CreateChapter(ctx context.Context, ch *domain.Chapter) error {
	query := `INSERT INTO chapters (id, story_id, chapter_number, title, content, illustrations, created_at, updated_at, has_content)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	if r.db.Driver == "sqlite" {
		query = `INSERT INTO chapters (id, story_id, chapter_number, title, content, illustrations, created_at, updated_at, has_content)
		         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := r.db.ExecContext(ctx, query,
		ch.ID,
		ch.StoryID,
		ch.ChapterNumber,
		ch.Title,
		ch.Content,
		ch.Illustrations,
		ch.CreatedAt,
		ch.UpdatedAt,
		contentFlag(ch.Content),
	)
	return err
}

func (r *StoryRepo) GetChapter(ctx context.Context, id string) (*domain.Chapter, error) {
	query := `SELECT id, story_id, chapter_number, title, content, illustrations, created_at, updated_at
	          FROM chapters WHERE id = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, story_id, chapter_number, title, content, illustrations, created_at, updated_at
		         FROM chapters WHERE id = ?`
	}

	var ch domain.Chapter
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ch.ID,
		&ch.StoryID,
		&ch.ChapterNumber,
		&ch.Title,
		&ch.Content,
		&ch.Illustrations,
		&ch.CreatedAt,
		&ch.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrChapterNotFound
		}
		return nil, err
	}

	return &ch, nil
}

func (r *StoryRepo) ListChaptersByStory(ctx context.Context, storyID string) ([]*domain.Chapter, error) {
	query := `SELECT id, story_id, chapter_number, title, content, illustrations, created_at, updated_at
	          FROM chapters WHERE story_id = $1 ORDER BY chapter_number ASC, created_at ASC`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, story_id, chapter_number, title, content, illustrations, created_at, updated_at
		         FROM chapters WHERE story_id = ? ORDER BY chapter_number ASC, created_at ASC`
	}

	rows, err := r.db.QueryContext(ctx, query, storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []*domain.Chapter
	for rows.Next() {
		var ch domain.Chapter
		if err := rows.Scan(
			&ch.ID,
			&ch.StoryID,
			&ch.ChapterNumber,
			&ch.Title,
			&ch.Content,
			&ch.Illustrations,
			&ch.CreatedAt,
			&ch.UpdatedAt,
		); err != nil {
			return nil, err
		}
		chapters = append(chapters, &ch)
	}

	return chapters, rows.Err()
}

func (r *StoryRepo) UpdateChapter(ctx context.Context, ch *domain.Chapter) error {
	query := `UPDATE chapters
	          SET title = $1, content = $2, illustrations = $3, updated_at = $4, has_content = $5
	          WHERE id = $6`
	if r.db.Driver == "sqlite" {
		query = `UPDATE chapters
		         SET title = ?, content = ?, illustrations = ?, updated_at = ?, has_content = ?
		         WHERE id = ?`
	}

	res, err := r.db.ExecContext(ctx, query,
		ch.Title,
		ch.Content,
		ch.Illustrations,
		ch.UpdatedAt,
		contentFlag(ch.Content),
		ch.ID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrChapterNotFound
	}
	return nil
}

func (r *StoryRepo) DeleteChapter(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, "DELETE FROM reading_progress WHERE chapter_id = $1", id); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, "DELETE FROM chapters WHERE id = $1", id)
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
	return tx.Commit()
}
