package storage

import (
	"database/sql"
	"os"
)

// migrate reads the SQL migration file and executes it against the database.
// Idempotent: uses CREATE TABLE IF NOT EXISTS and CREATE INDEX IF NOT EXISTS,
// so running it multiple times is safe.
func migrate(db *sql.DB) error {
	data, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))
	return err
}
