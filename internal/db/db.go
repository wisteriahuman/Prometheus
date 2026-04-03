package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) *DB {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	os.MkdirAll(dir, 0o755)

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic("failed to open database: " + err.Error())
	}

	// Enable WAL mode and foreign keys
	conn.Exec("PRAGMA journal_mode = WAL")
	conn.Exec("PRAGMA foreign_keys = ON")

	return &DB{conn}
}
