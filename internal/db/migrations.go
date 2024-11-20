package db

import (
	"database/sql"
	"os"
	"path/filepath"
)

func Migrate(db *sql.DB) error {
	var err error;

	path := filepath.Join("internal", "db", "create_tables.sql")
	c, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	query := string(c)

	if _, err = db.Exec(query); err != nil {
		return err
	}

	return err
}