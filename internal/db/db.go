package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type DBHandler struct {
	DB *sql.DB
}

// Connect to DB instance
func Connect(dbFile string) (*sql.DB, error) {
	var err error;
	
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	return db, err
}