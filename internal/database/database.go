package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const DATABASE_PATH = "storage/database.sqlite"

var DB *sql.DB

func InitDatabase() *sql.DB {
	if DB == nil {
		if err := os.MkdirAll(filepath.Dir(DATABASE_PATH), 0o755); err != nil {
			panic("cannot create database directory")
		}

		db, err := sql.Open("sqlite", DATABASE_PATH)

		if err != nil {
			panic("no db")
		}

		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
			)
		`)

		if err != nil {
			db.Close()
			panic("cannot create users table")
		}

		DB = db
	}

	return DB
}
