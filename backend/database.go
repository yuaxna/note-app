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
	seedData()

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

func seedData() {
	// Replace these hashed passwords with real bcrypt hashes for login
	const user1Password = "$2a$10$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" // dummy hash
	const user2Password = "$2a$10$YYYYYYYYYYYYYYYYYYYYYYYYYYYYYY" // dummy hash

	_, err := DB.Exec(`
		INSERT OR IGNORE INTO users (id, fullname, email, username, password, gender)
		VALUES
			(1, 'User One', 'user1@example.com', 'user1', ?, 'Female'),
			(2, 'User Two', 'user2@example.com', 'user2', ?, 'Male')
	`, user1Password, user2Password)
	if err != nil {
		log.Fatal("Failed to insert seed users:", err)
	}

	_, err = DB.Exec(`
		INSERT OR IGNORE INTO notes (id, title, content, user_id)
		VALUES
			(1, 'First Note', 'This is the first note content', 1),
			(2, 'Second Note', 'This is the second note content', 2)
	`)
	if err != nil {
		log.Fatal("Failed to insert seed notes:", err)
	}

	_, err = DB.Exec(`
		INSERT OR IGNORE INTO shared_notes (note_id, user_id, can_edit)
		VALUES
			(2, 1, 1)
	`)
	if err != nil {
		log.Fatal("Failed to insert seed shared_notes:", err)
	}
}
