package backend

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader config to upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins or you can restrict this for security
		return true
	},
}

func WSHandler(c *gin.Context) {
	// Check user authentication (must be logged in)
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Upgrade to WebSocket *only if authorized*
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer conn.Close()

	// Handle messages (echo or broadcast later)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break // connection closed or error
		}

		// For now, just echo back the message
		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
