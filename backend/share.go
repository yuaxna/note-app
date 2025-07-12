package backend

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShareNote(c *gin.Context) {
	var input ShareInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	currentUserID := c.GetInt("user_id")

	// Confirm note belongs to user
	row := DB.QueryRow("SELECT id FROM notes WHERE id = ? AND user_id = ?", input.NoteID, currentUserID)
	var noteID int
	if err := row.Scan(&noteID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Note not found or not owned by you"})
		return
	}

	// Prevent sharing with self
	if input.TargetUserID == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot share with yourself"})
		return
	}

	// <-- Add the log here -->
	log.Println("Sharing note:", input.NoteID, "to user:", input.TargetUserID)

	// Insert into shared_notes (ignoring duplicates due to UNIQUE constraint)
	_, err := DB.Exec(`
        INSERT OR IGNORE INTO shared_notes (note_id, user_id) VALUES (?, ?)
    `, input.NoteID, input.TargetUserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note shared successfully"})
}

func GetSharedNotes(c *gin.Context) {
	userID := c.GetInt("user_id")
	log.Printf("Fetching shared notes for user_id: %d", userID)

	rows, err := DB.Query(`
        SELECT n.id, n.title, n.content, u.username, n.created_at
        FROM shared_notes s
        JOIN notes n ON s.note_id = n.id
        JOIN users u ON n.user_id = u.id
        WHERE s.user_id = ?
    `, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shared notes"})
		return
	}
	defer rows.Close()

	var shared []gin.H
	for rows.Next() {
		var id int
		var title, content, author, createdAt string
		if err := rows.Scan(&id, &title, &content, &author, &createdAt); err != nil {
			log.Println("Error scanning shared note:", err)
			continue
		}
		shared = append(shared, gin.H{
			"id":         id,
			"title":      title,
			"content":    content,
			"author":     author,
			"created_at": createdAt,
		})
	}

	if shared == nil {
		shared = []gin.H{}
	}

	c.JSON(http.StatusOK, shared)
}

func GetUsers(c *gin.Context) {
	currentUserID := c.GetInt("user_id")

	rows, err := DB.Query("SELECT id, fullname, username FROM users WHERE id != ?", currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var fullname, username string
		if err := rows.Scan(&id, &fullname, &username); err != nil {
			continue
		}
		users = append(users, map[string]interface{}{
			"id":       id,
			"fullname": fullname,
			"username": username,
		})
	}

	c.JSON(http.StatusOK, users)
}
