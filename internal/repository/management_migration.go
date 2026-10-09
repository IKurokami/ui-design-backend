package repository

func (db *DB) migrateManagement() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(64) PRIMARY KEY)`); err != nil {
		return err
	}
	var applied int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = 'management-v1'`).Scan(&applied); err != nil {
		return err
	}
	if applied != 0 {
		return tx.Commit()
	}
	if _, err := tx.Exec(`ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT '' CHECK (role IN ('', 'admin', 'moderator'))`); err != nil {
		return err
	}
	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`); err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES ('management-v1')`); err != nil {
		return err
	}
	return tx.Commit()
}
