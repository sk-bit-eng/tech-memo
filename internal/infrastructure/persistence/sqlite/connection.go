package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate sqlite database: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS memos (
			id          TEXT PRIMARY KEY,
			user_id     TEXT NOT NULL,
			title       TEXT NOT NULL,
			content     TEXT NOT NULL,
			category_id TEXT NOT NULL DEFAULT '',
			parameters  TEXT NOT NULL DEFAULT '[]',
			is_pinned   INTEGER NOT NULL DEFAULT 0,
			created_at  DATETIME NOT NULL,
			updated_at  DATETIME NOT NULL,
			deleted_at  DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS todos (
			id           TEXT PRIMARY KEY,
			user_id      TEXT NOT NULL,
			title        TEXT NOT NULL,
			content      TEXT NOT NULL,
			category_id  TEXT NOT NULL DEFAULT '',
			parameters   TEXT NOT NULL DEFAULT '[]',
			is_pinned    INTEGER NOT NULL DEFAULT 0,
			due_at       DATETIME,
			completed_at DATETIME,
			created_at   DATETIME NOT NULL,
			updated_at   DATETIME NOT NULL,
			deleted_at   DATETIME
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
