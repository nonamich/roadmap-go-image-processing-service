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

		db, err := sql.Open("sqlite", DATABASE_PATH+"?_pragma=foreign_keys(1)")

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

		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS images (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				uuid TEXT NOT NULL UNIQUE,
				metadata JSON NOT NULL,
				user_id INTEGER NOT NULL,
				CONSTRAINT fk_images_user
					FOREIGN KEY (user_id) REFERENCES users(id)
					ON DELETE CASCADE
			)
		`)

		if err != nil {
			db.Close()
			panic("cannot create images table")
		}

		DB = db
	}

	return DB
}
