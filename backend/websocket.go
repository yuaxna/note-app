package backend

import (
	"time"

	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust origin check if needed for security
	},
}

// NoteUpdateMessage struct - message to send to clients

var Manager = ClientManager{
	clients:    make(map[*websocket.Conn]bool),
	broadcast:  make(chan NoteUpdateMessage),
	register:   make(chan *websocket.Conn),
	unregister: make(chan *websocket.Conn),
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.mu.Lock()
			manager.clients[conn] = true
			manager.mu.Unlock()
		case conn := <-manager.unregister:
			manager.mu.Lock()
			if _, ok := manager.clients[conn]; ok {
				delete(manager.clients, conn)
				conn.Close()
			}
			manager.mu.Unlock()
		case message := <-manager.broadcast:
			manager.mu.Lock()
			for conn := range manager.clients {
				err := conn.WriteJSON(message)
				if err != nil {
					conn.Close()
					delete(manager.clients, conn)
				}
			}
			manager.mu.Unlock()
		}
	}
}

// WSHandler upgrades HTTP connection to websocket and manages connection lifecycle
func WSHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	username := session.Get("username") // make sure you store username on login

	if userID == nil || username == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}

	Manager.register <- conn

	defer func() {
		Manager.unregister <- conn
	}()

	// Just read messages (if needed) or ignore incoming messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// Call this function in your note create/update/delete/share handlers
func BroadcastNoteUpdate(action string, noteID int, title, content, sender string) {
	Manager.broadcast <- NoteUpdateMessage{
		Action:    action,
		NoteID:    noteID,
		Title:     title,
		Content:   content,
		Sender:    sender,
		Timestamp: time.Now(),
	}
}

// In backend/websocket.go
func StartManager() {
    go Manager.start()
}
