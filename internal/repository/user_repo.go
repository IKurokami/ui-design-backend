package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	query := `INSERT INTO users (id, display_name, email, password_hash, created_at, updated_at, role)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if r.db.Driver == "sqlite" {
		query = `INSERT INTO users (id, display_name, email, password_hash, created_at, updated_at, role)
		         VALUES (?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := r.db.ExecContext(ctx, query, u.ID, u.DisplayName, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt, u.Role)
	return err
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, display_name, email, password_hash, created_at, updated_at, role FROM users WHERE id = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, display_name, email, password_hash, created_at, updated_at, role FROM users WHERE id = ?`
	}

	row := r.db.QueryRowContext(ctx, query, id)
	var u domain.User
	err := row.Scan(&u.ID, &u.DisplayName, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, display_name, email, password_hash, created_at, updated_at, role FROM users WHERE email = $1`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, display_name, email, password_hash, created_at, updated_at, role FROM users WHERE email = ?`
	}

	row := r.db.QueryRowContext(ctx, query, email)
	var u domain.User
	err := row.Scan(&u.ID, &u.DisplayName, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, display_name, email, password_hash, created_at, updated_at, role FROM users ORDER BY created_at DESC, id ASC`
	if r.db.Driver == "sqlite" {
		query = `SELECT id, display_name, email, password_hash, created_at, updated_at, role FROM users ORDER BY created_at DESC, id ASC`
	}
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.DisplayName, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}

func (r *UserRepo) ListPage(ctx context.Context, page, pageSize int, query, role string) ([]*domain.User, int, error) {
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
	where := " WHERE (display_name LIKE $1 OR email LIKE $1)"
	args := []any{like}
	if r.db.Driver == "sqlite" {
		where = " WHERE (display_name LIKE ? OR email LIKE ?)"
		args = []any{like, like}
	}
	if role != "" {
		where += " AND role = " + placeholder(r.db.Driver, len(args)+1)
		args = append(args, role)
	}
	countQuery := "SELECT COUNT(*) FROM users" + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	limitP, offsetP := placeholder(r.db.Driver, len(args)+1), placeholder(r.db.Driver, len(args)+2)
	q := "SELECT id, display_name, email, password_hash, created_at, updated_at, role FROM users" + where + " ORDER BY created_at DESC, id ASC LIMIT " + limitP + " OFFSET " + offsetP
	args = append(args, pageSize, offset)
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]*domain.User, 0, pageSize)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.DisplayName, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt, &u.Role); err != nil {
			return nil, 0, err
		}
		items = append(items, &u)
	}
	return items, total, rows.Err()
}

func placeholder(driver string, n int) string {
	if driver == "sqlite" {
		return "?"
	}
	return "$" + fmt.Sprint(n)
}

func (r *UserRepo) UpdateRole(ctx context.Context, userID string, role domain.Role) error {
	query := `UPDATE users SET role = $1, updated_at = $2 WHERE id = $3`
	if r.db.Driver == "sqlite" {
		query = `UPDATE users SET role = ?, updated_at = ? WHERE id = ?`
	}
	res, err := r.db.ExecContext(ctx, query, role, time.Now().UTC(), userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var exists string
	userQuery := `SELECT id FROM users WHERE id = $1`
	if r.db.Driver == "sqlite" {
		userQuery = `SELECT id FROM users WHERE id = ?`
	}
	if err := tx.QueryRowContext(ctx, userQuery, userID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	placeholder := "$1"
	if r.db.Driver == "sqlite" {
		placeholder = "?"
	}
	queries := []string{
		"DELETE FROM reading_progress WHERE user_id = " + placeholder,
		"DELETE FROM reading_progress WHERE story_id IN (SELECT id FROM stories WHERE user_id = " + placeholder + ")",
		"DELETE FROM chapters WHERE story_id IN (SELECT id FROM stories WHERE user_id = " + placeholder + ")",
		"DELETE FROM stories WHERE user_id = " + placeholder,
		"DELETE FROM ebooks WHERE user_id = " + placeholder,
		"DELETE FROM user_profiles WHERE user_id = " + placeholder,
		"DELETE FROM password_resets WHERE user_id = " + placeholder,
		"DELETE FROM sessions WHERE user_id = " + placeholder,
		"DELETE FROM users WHERE id = " + placeholder,
	}
	for _, query := range queries {
		if _, err := tx.ExecContext(ctx, query, userID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID, newPasswordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	if r.db.Driver == "sqlite" {
		query = `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	}

	_, err := r.db.ExecContext(ctx, query, newPasswordHash, time.Now(), userID)
	return err
}

func (r *UserRepo) UpdateDisplayName(ctx context.Context, userID, displayName string) error {
	query := `UPDATE users SET display_name = $1, updated_at = $2 WHERE id = $3`
	if r.db.Driver == "sqlite" {
		query = `UPDATE users SET display_name = ?, updated_at = ? WHERE id = ?`
	}

	_, err := r.db.ExecContext(ctx, query, displayName, time.Now().UTC(), userID)
	return err
}
