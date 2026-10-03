package storage

import (
	"database/sql"

	"github.com/Carlos20052030/cloudmonitor/migrations"
)

// migrate applies the initial schema to the database.
// The SQL is embedded in the binary, so it works regardless of the
// process working directory.
func migrate(db *sql.DB) error {
	data, err := migrations.FS.ReadFile("001_init.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))
	return err
}
