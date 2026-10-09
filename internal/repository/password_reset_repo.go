package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ui-design-backend/internal/domain"
)

var ErrResetNotFound = errors.New("password reset token not found")

type PasswordResetRepo struct {
	db *DB
}

func NewPasswordResetRepo(db *DB) *PasswordResetRepo {
	return &PasswordResetRepo{db: db}
}

func (r *PasswordResetRepo) Create(ctx context.Context, pr *domain.PasswordReset) error {
	query := `INSERT INTO password_resets (id, user_id, token_hash, expires_at, created_at, used_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	if r.db.Driver == "sqlite" {
		query = `INSERT INTO password_resets (id, user_id, token_hash, expires_at, created_at, used_at) 
		         VALUES (?, ?, ?, ?, ?, ?)`
	}

	_, err := r.db.ExecContext(ctx, query, pr.ID, pr.UserID, pr.TokenHash, pr.ExpiresAt, pr.CreatedAt, pr.UsedAt)
	return err
}

func (r *PasswordResetRepo) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordReset, error) {
	query := `SELECT id, user_id, token_hash, expires_at, created_at, used_at FROM password_resets WHERE token_hash = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, user_id, token_hash, expires_at, created_at, used_at FROM password_resets WHERE token_hash = ?`
	}

	row := r.db.QueryRowContext(ctx, query, tokenHash)
	var pr domain.PasswordReset
	var usedAt sql.NullTime

	err := row.Scan(&pr.ID, &pr.UserID, &pr.TokenHash, &pr.ExpiresAt, &pr.CreatedAt, &usedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrResetNotFound
		}
		return nil, err
	}

	if usedAt.Valid {
		pr.UsedAt = &usedAt.Time
	}
	return &pr, nil
}

// Complete consumes the token, updates the password and revokes sessions atomically.
func (r *PasswordResetRepo) Complete(ctx context.Context, id, passwordHash string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	query := `UPDATE password_resets SET used_at = $1 WHERE id = $2 AND used_at IS NULL RETURNING user_id, expires_at`
	if r.db.Driver == "sqlite" {
		query = `UPDATE password_resets SET used_at = ? WHERE id = ? AND used_at IS NULL RETURNING user_id, expires_at`
	}
	var userID string
	var expiresAt time.Time
	if err := tx.QueryRowContext(ctx, query, now, id).Scan(&userID, &expiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrResetNotFound
		}
		return err
	}
	// Check after acquiring the write lock, including time spent waiting for it.
	if !time.Now().Before(expiresAt) {
		return ErrResetNotFound
	}

	query = `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	if r.db.Driver == "sqlite" {
		query = `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	}
	result, err := tx.ExecContext(ctx, query, passwordHash, now, userID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrUserNotFound
	}

	query = `UPDATE sessions SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	if r.db.Driver == "sqlite" {
		query = `UPDATE sessions SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`
	}
	if _, err := tx.ExecContext(ctx, query, now, userID); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *PasswordResetRepo) InvalidateActiveForUser(ctx context.Context, userID string) error {
	query := `UPDATE password_resets SET used_at = $1 WHERE user_id = $2 AND used_at IS NULL`
	if r.db.Driver == "sqlite" {
		query = `UPDATE password_resets SET used_at = ? WHERE user_id = ? AND used_at IS NULL`
	}

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}
