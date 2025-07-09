package backend

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // or whatever database driver you're using
)

var DB *sql.DB

func InitDB() {
	var err error
	// Using SQLite for this example - change connection string for your database
	DB, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create tables if they don't exist
	createTables()
	log.Println("Database connected successfully")
}

func createTables() {
	// Create users table
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

	// Create notes table
	noteTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);`

	if _, err := DB.Exec(userTable); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	if _, err := DB.Exec(noteTable); err != nil {
		log.Fatal("Failed to create notes table:", err)
	}

	log.Println("Database tables created/verified successfully")
}
