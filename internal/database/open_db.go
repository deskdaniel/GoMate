package database

import (
	"database/sql"
	"fmt"
)

func OpenDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./chess.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys in db: %w", err)
	}

	return db, nil
}
