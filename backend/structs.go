package backend

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	ID         int       `json:"id"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Gender     string    `json:"gender"`
	Identifier string    `json:"identifier"` // For login (username or email)
	CreatedAt  time.Time `json:"created_at"`
}

type LoginInput struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Note struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username"`
}

type ShareInput struct {
	NoteID       int `json:"note_id"`
	TargetUserID int `json:"target_user_id"`
}

type NoteUpdateMessage struct {
	Action    string    `json:"action"`
	NoteID    int       `json:"note_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
}

// ClientManager manages WebSocket connections
type ClientManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan NoteUpdateMessage
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

type NoteUpdateInput struct {
	ID      int    `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
