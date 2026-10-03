package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ui-design-backend/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (id, display_name, email, password_hash, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	if r.db.Driver == "sqlite" {
		query = `INSERT INTO users (id, display_name, email, password_hash, created_at, updated_at) 
		         VALUES (?, ?, ?, ?, ?, ?)`
	}

	_, err := r.db.ExecContext(ctx, query, u.ID, u.DisplayName, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, display_name, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, display_name, email, password_hash, created_at, updated_at FROM users WHERE id = ?`
	}

	row := r.db.QueryRowContext(ctx, query, id)
	var u domain.User
	err := row.Scan(&u.ID, &u.DisplayName, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, display_name, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, display_name, email, password_hash, created_at, updated_at FROM users WHERE email = ?`
	}

	row := r.db.QueryRowContext(ctx, query, email)
	var u domain.User
	err := row.Scan(&u.ID, &u.DisplayName, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID, newPasswordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	if r.db.Driver == "sqlite" {
		query = `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	}

	_, err := r.db.ExecContext(ctx, query, newPasswordHash, time.Now(), userID)
	return err
}
