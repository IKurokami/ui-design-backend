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

	CREATE TABLE IF NOT EXISTS user_profiles (
		id VARCHAR(64) PRIMARY KEY,
		user_id VARCHAR(64) NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
		avatar_url TEXT,
		banner_url TEXT,
		bio TEXT,
		headline VARCHAR(150),
		theme_color VARCHAR(30) DEFAULT '#3B82F6',
		custom_css TEXT,
		website_url VARCHAR(255),
		location VARCHAR(100),
		social_links TEXT,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);

	CREATE TABLE IF NOT EXISTS stories (
		id VARCHAR(64) PRIMARY KEY,
		user_id VARCHAR(64) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		cover_url TEXT,
		category VARCHAR(100),
		tags TEXT,
		status VARCHAR(50) DEFAULT 'published',
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_stories_user_id ON stories(user_id);

	CREATE TABLE IF NOT EXISTS chapters (
		id VARCHAR(64) PRIMARY KEY,
		story_id VARCHAR(64) NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
		chapter_number INT NOT NULL DEFAULT 1,
		title VARCHAR(255) NOT NULL,
		content TEXT,
		illustrations TEXT,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_chapters_story_id ON chapters(story_id);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	if err := db.migrateManagement(); err != nil {
		return err
	}
	if err := db.migrateReading(); err != nil {
		return err
	}
	if err := db.migrateChapterOnly(); err != nil {
		return err
	}
	if err := db.migrateEbooks(); err != nil { return err }
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS favorite_stories (
		user_id VARCHAR(64) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		story_id VARCHAR(64) NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL,
		PRIMARY KEY (user_id, story_id)
	);
	CREATE INDEX IF NOT EXISTS idx_favorite_stories_user_time ON favorite_stories(user_id, created_at DESC);`)
	return err
}
