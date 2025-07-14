package backend

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
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
	note.CreatedAt = time.Now().Format("2006-01-02 15:04:05") // Use a more standard format
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

func GetUserNotes(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := DB.Query(`
    SELECT 
        n.id, n.title, n.content,
        COALESCE(n.created_at, datetime('now')) AS created_at,
        u.username,
        CASE WHEN n.user_id = ? THEN 1 ELSE 0 END AS is_owner,
        CASE 
  			WHEN n.user_id = ? THEN 1
  			WHEN s.user_id = ? THEN 1
  			ELSE 0
		END AS can_edit
    FROM notes n
    LEFT JOIN shared_notes s ON s.note_id = n.id AND s.user_id = ?
    JOIN users u ON n.user_id = u.id
    WHERE n.user_id = ? OR s.user_id = ?
    ORDER BY n.created_at DESC
`, userID, userID, userID, userID, userID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}
	defer rows.Close()

	var notes []map[string]interface{}
	for rows.Next() {
		var (
			id        int
			title     string
			content   string
			createdAt string
			username  string
			isOwner   int
			canEdit   int
		)

		if err := rows.Scan(&id, &title, &content, &createdAt, &username, &isOwner, &canEdit); err != nil {
			continue
		}

		note := map[string]interface{}{
			"id":         id,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"username":   username,
			"is_owner":   isOwner == 1,
			"can_edit":   canEdit == 1,
		}
		notes = append(notes, note)
	}

	c.JSON(http.StatusOK, notes)
}

func CreateNote(c *gin.Context) {
	var note Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetInt("user_id")

	_, err := DB.Exec(
		"INSERT INTO notes (title, content, user_id) VALUES (?, ?, ?)",
		note.Title, note.Content, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note created successfully"})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			if IsAPIRequest(c) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			} else {
				c.Redirect(http.StatusFound, "/")
			}
			c.Abort()
			return
		}
		if id, ok := userID.(int); ok {
			c.Set("user_id", id)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

func UpdateNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	currentUserID := c.GetInt("user_id")
	log.Printf("[DEBUG] Checking edit permission for user %d on note %d", currentUserID, id)

	// --- Check permission to edit note ---
	var exists bool
	err = DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM notes n
			LEFT JOIN shared_notes s ON s.note_id = n.id AND s.user_id = ?
			WHERE n.id = ? AND (n.user_id = ? OR s.user_id = ?)
		)
	`, currentUserID, id, currentUserID, currentUserID).Scan(&exists)

	log.Printf("[DEBUG] Permission check error: %v", err)
	log.Printf("[DEBUG] Permission exists: %v", exists)

	if err != nil {
		log.Printf("[ERROR] DB error during permission check: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during permission check"})
		return
	}

	log.Printf("[DEBUG] Permission exists: %v", exists)

	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to edit this note"})
		return
	}

	// --- Update note ---
	_, err = DB.Exec(`
        UPDATE notes SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
    `, input.Title, input.Content, id)

	if err != nil {
		log.Printf("[ERROR] Failed to update note: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

func DeleteNote(c *gin.Context) {
	// Get note ID from URL param
	noteID := c.Param("id")
	userID := c.GetInt("user_id")

	res, err := DB.Exec("DELETE FROM notes WHERE id = ? AND user_id = ?", noteID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Note not found or not owned by user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
