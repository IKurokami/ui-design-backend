package repository

import (
	"context"
	"database/sql"
	"errors"

	"ui-design-backend/internal/domain"
)

var ErrEbookNotFound = errors.New("ebook not found")
var ErrEbookConflict = errors.New("ebook revision conflict")

func (db *DB) migrateEbooks() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ebooks (
 id VARCHAR(64) PRIMARY KEY,
 user_id VARCHAR(64) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 title VARCHAR(255) NOT NULL,
 cover_url TEXT NOT NULL DEFAULT '',
 page_count INT NOT NULL,
 revision INT NOT NULL DEFAULT 1,
 design TEXT NOT NULL,
 created_at TIMESTAMP NOT NULL,
 updated_at TIMESTAMP NOT NULL
 ); CREATE INDEX IF NOT EXISTS idx_ebooks_user_id ON ebooks(user_id);`)
	return err
}

type EbookRepo struct{ db *DB }

func NewEbookRepo(db *DB) *EbookRepo { return &EbookRepo{db: db} }

func (r *EbookRepo) query(sqlite, postgres string) string {
	if r.db.Driver == "sqlite" {
		return sqlite
	}
	return postgres
}

func (r *EbookRepo) Create(ctx context.Context, e *domain.Ebook) error {
	_, err := r.db.ExecContext(ctx, r.query(
		`INSERT INTO ebooks (id, user_id, title, cover_url, page_count, revision, design, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?)`,
		`INSERT INTO ebooks (id, user_id, title, cover_url, page_count, revision, design, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
	), e.ID, e.UserID, e.Title, e.CoverURL, e.PageCount, e.Revision, string(e.Design), e.CreatedAt, e.UpdatedAt)
	return err
}

func (r *EbookRepo) Get(ctx context.Context, userID, id string) (*domain.Ebook, error) {
	var e domain.Ebook
	var design string
	err := r.db.QueryRowContext(ctx, r.query(
		`SELECT id, user_id, title, cover_url, page_count, revision, design, created_at, updated_at FROM ebooks WHERE id=? AND user_id=?`,
		`SELECT id, user_id, title, cover_url, page_count, revision, design, created_at, updated_at FROM ebooks WHERE id=$1 AND user_id=$2`,
	), id, userID).Scan(&e.ID, &e.UserID, &e.Title, &e.CoverURL, &e.PageCount, &e.Revision, &design, &e.CreatedAt, &e.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEbookNotFound
	}
	if err != nil {
		return nil, err
	}
	e.Design = []byte(design)
	return &e, nil
}

func (r *EbookRepo) List(ctx context.Context, userID string) ([]*domain.Ebook, error) {
	rows, err := r.db.QueryContext(ctx, r.query(
		`SELECT id, user_id, title, cover_url, page_count, revision, created_at, updated_at FROM ebooks WHERE user_id=? ORDER BY updated_at DESC, id ASC`,
		`SELECT id, user_id, title, cover_url, page_count, revision, created_at, updated_at FROM ebooks WHERE user_id=$1 ORDER BY updated_at DESC, id ASC`,
	), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ebooks := make([]*domain.Ebook, 0)
	for rows.Next() {
		var e domain.Ebook
		if err := rows.Scan(&e.ID, &e.UserID, &e.Title, &e.CoverURL, &e.PageCount, &e.Revision, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		ebooks = append(ebooks, &e)
	}
	return ebooks, rows.Err()
}

func (r *EbookRepo) Update(ctx context.Context, e *domain.Ebook, expectedRevision int) error {
	result, err := r.db.ExecContext(ctx, r.query(
		`UPDATE ebooks SET title=?, cover_url=?, page_count=?, design=?, revision=revision+1, updated_at=? WHERE id=? AND user_id=? AND revision=?`,
		`UPDATE ebooks SET title=$1, cover_url=$2, page_count=$3, design=$4, revision=revision+1, updated_at=$5 WHERE id=$6 AND user_id=$7 AND revision=$8`,
	), e.Title, e.CoverURL, e.PageCount, string(e.Design), e.UpdatedAt, e.ID, e.UserID, expectedRevision)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrEbookConflict
	}
	e.Revision = expectedRevision + 1
	return nil
}

func (r *EbookRepo) Delete(ctx context.Context, userID, id string) error {
	result, err := r.db.ExecContext(ctx, r.query(`DELETE FROM ebooks WHERE id=? AND user_id=?`, `DELETE FROM ebooks WHERE id=$1 AND user_id=$2`), id, userID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrEbookNotFound
	}
	return nil
}
