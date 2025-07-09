package backend

import (
	"net/http"

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

// WSHandler upgrades the HTTP connection to WebSocket and handles messages
func WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer conn.Close()

	// Now handle incoming messages and broadcast or save changes as needed
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break // connection closed or error
		}

		// Here you can parse msg and do stuff with it, e.g. update notes live

		// For now, echo back message (test)
		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
