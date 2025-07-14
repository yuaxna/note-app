package backend

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTables()

	log.Println("Database connected, tables created/verified, and seed data inserted successfully")
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fullname TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		gender TEXT DEFAULT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	noteTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);`

	sharedNotesTable := `
	CREATE TABLE IF NOT EXISTS shared_notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		note_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		can_edit INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (note_id, user_id)
	);`

	if _, err := DB.Exec(userTable); err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	if _, err := DB.Exec(noteTable); err != nil {
		log.Fatal("Failed to create notes table:", err)
	}
	if _, err := DB.Exec(sharedNotesTable); err != nil {
		log.Fatal("Failed to create shared_notes table:", err)
	}
}
