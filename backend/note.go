package backend

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// --- CREATE a note ---
func CreateNoteHandler(c *gin.Context) {
	var note Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetInt("userID") // Set by auth middleware
	note.UserID = userID
	note.CreatedAt = time.Now().Format(time.RFC3339)
	note.UpdatedAt = note.CreatedAt

	stmt, err := DB.Prepare("INSERT INTO notes(user_id, title, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(note.UserID, note.Title, note.Content, note.CreatedAt, note.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Note created"})
}
