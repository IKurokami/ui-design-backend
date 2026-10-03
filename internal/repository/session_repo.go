package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ui-design-backend/internal/domain"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionRepo struct {
	db *DB
}

func NewSessionRepo(db *DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(ctx context.Context, s *domain.Session) error {
	query := `INSERT INTO sessions (id, user_id, token_hash, user_agent, ip, expires_at, created_at, revoked_at, rotated_to_id, grace_period_expires_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	if r.db.Driver == "sqlite" {
		query = `INSERT INTO sessions (id, user_id, token_hash, user_agent, ip, expires_at, created_at, revoked_at, rotated_to_id, grace_period_expires_at) 
		         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := r.db.ExecContext(ctx, query, s.ID, s.UserID, s.TokenHash, s.UserAgent, s.IP, s.ExpiresAt, s.CreatedAt, s.RevokedAt, s.RotatedToID, s.GracePeriodExpiresAt)
	return err
}

func (r *SessionRepo) FindByID(ctx context.Context, id string) (*domain.Session, error) {
	query := `SELECT id, user_id, token_hash, user_agent, ip, expires_at, created_at, revoked_at, rotated_to_id, grace_period_expires_at 
	          FROM sessions WHERE id = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, user_id, token_hash, user_agent, ip, expires_at, created_at, revoked_at, rotated_to_id, grace_period_expires_at 
		         FROM sessions WHERE id = ?`
	}

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanSession(row)
}

func (r *SessionRepo) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	query := `SELECT id, user_id, token_hash, user_agent, ip, expires_at, created_at, revoked_at, rotated_to_id, grace_period_expires_at 
	          FROM sessions WHERE token_hash = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, user_id, token_hash, user_agent, ip, expires_at, created_at, revoked_at, rotated_to_id, grace_period_expires_at 
		         FROM sessions WHERE token_hash = ?`
	}

	row := r.db.QueryRowContext(ctx, query, tokenHash)
	return r.scanSession(row)
}

func (r *SessionRepo) Rotate(ctx context.Context, oldSessionID string, newSession *domain.Session, gracePeriod time.Duration) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	graceUntil := now.Add(gracePeriod)

	// 1. Mark old session as rotated with grace period
	updateQuery := `UPDATE sessions SET rotated_to_id = $1, grace_period_expires_at = $2 WHERE id = $3`
	if r.db.Driver == "sqlite" {
		updateQuery = `UPDATE sessions SET rotated_to_id = ?, grace_period_expires_at = ? WHERE id = ?`
	}
	if _, err := tx.ExecContext(ctx, updateQuery, newSession.ID, graceUntil, oldSessionID); err != nil {
		return err
	}

	// 2. Insert new session
	insertQuery := `INSERT INTO sessions (id, user_id, token_hash, user_agent, ip, expires_at, created_at, revoked_at, rotated_to_id, grace_period_expires_at) 
	                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	if r.db.Driver == "sqlite" {
		insertQuery = `INSERT INTO sessions (id, user_id, token_hash, user_agent, ip, expires_at, created_at, revoked_at, rotated_to_id, grace_period_expires_at) 
		               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}
	if _, err := tx.ExecContext(ctx, insertQuery, newSession.ID, newSession.UserID, newSession.TokenHash, newSession.UserAgent, newSession.IP, newSession.ExpiresAt, newSession.CreatedAt, newSession.RevokedAt, newSession.RotatedToID, newSession.GracePeriodExpiresAt); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SessionRepo) Revoke(ctx context.Context, id string) error {
	query := `UPDATE sessions SET revoked_at = $1 WHERE id = $2`
	if r.db.Driver == "sqlite" {
		query = `UPDATE sessions SET revoked_at = ? WHERE id = ?`
	}
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *SessionRepo) RevokeAllForUser(ctx context.Context, userID string) error {
	query := `UPDATE sessions SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	if r.db.Driver == "sqlite" {
		query = `UPDATE sessions SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`
	}
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

func (r *SessionRepo) scanSession(row *sql.Row) (*domain.Session, error) {
	var s domain.Session
	var revokedAt sql.NullTime
	var rotatedToID sql.NullString
	var gracePeriodExpiresAt sql.NullTime

	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.TokenHash,
		&s.UserAgent,
		&s.IP,
		&s.ExpiresAt,
		&s.CreatedAt,
		&revokedAt,
		&rotatedToID,
		&gracePeriodExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if revokedAt.Valid {
		s.RevokedAt = &revokedAt.Time
	}
	if rotatedToID.Valid {
		s.RotatedToID = &rotatedToID.String
	}
	if gracePeriodExpiresAt.Valid {
		s.GracePeriodExpiresAt = &gracePeriodExpiresAt.Time
	}

	return &s, nil
}
