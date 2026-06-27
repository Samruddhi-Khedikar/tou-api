package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// New opens (or creates) the SQLite database file and applies the schema.
func New() (*sql.DB, error) {
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "tou.db"
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// SQLite is not great with concurrent writers; one connection is enough for a service this size.
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if err := applySchema(db); err != nil {
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	return db, nil
}

func applySchema(db *sql.DB) error {
	// Validations are handled at the handler level via go-playground/validator.
	// Schema only enforces structural constraints (NOT NULL, foreign keys).
	schema := `
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS chargers (
		id          TEXT PRIMARY KEY,
		name        TEXT NOT NULL,
		location    TEXT NOT NULL,
		timezone    TEXT NOT NULL DEFAULT 'UTC',
		created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
		updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
	);

	CREATE TABLE IF NOT EXISTS pricing_periods (
		id            TEXT PRIMARY KEY,
		charger_id    TEXT NOT NULL REFERENCES chargers(id) ON DELETE CASCADE,
		start_hour    INTEGER NOT NULL,
		start_minute  INTEGER NOT NULL DEFAULT 0,
		end_hour      INTEGER NOT NULL,
		end_minute    INTEGER NOT NULL DEFAULT 0,
		price_per_kwh REAL NOT NULL,
		created_at    TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
	);

	CREATE INDEX IF NOT EXISTS idx_pricing_periods_charger ON pricing_periods(charger_id);
	`
	_, err := db.Exec(schema)
	return err
}
