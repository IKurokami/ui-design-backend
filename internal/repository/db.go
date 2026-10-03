package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
	Driver string
}

func NewDB(driver, dsn string) (*DB, error) {
	if driver == "sqlite" {
		dir := filepath.Dir(dsn)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create sqlite dir: %w", err)
			}
		}
	}

	dbDriver := driver
	if driver == "postgres" {
		dbDriver = "pgx"
	}

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	res := &DB{DB: db, Driver: driver}
	if err := res.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return res, nil
}

func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(64) PRIMARY KEY,
		display_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id VARCHAR(64) PRIMARY KEY,
		user_id VARCHAR(64) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token_hash VARCHAR(64) NOT NULL UNIQUE,
		user_agent TEXT,
		ip VARCHAR(45),
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL,
		revoked_at TIMESTAMP,
		rotated_to_id VARCHAR(64),
		grace_period_expires_at TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON sessions(token_hash);
	CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

	CREATE TABLE IF NOT EXISTS password_resets (
		id VARCHAR(64) PRIMARY KEY,
		user_id VARCHAR(64) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token_hash VARCHAR(64) NOT NULL UNIQUE,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL,
		used_at TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_password_resets_token_hash ON password_resets(token_hash);
	CREATE INDEX IF NOT EXISTS idx_password_resets_user_id ON password_resets(user_id);
	`

	_, err := db.Exec(schema)
	return err
}
