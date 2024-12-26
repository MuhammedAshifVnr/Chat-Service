package models

import (
	"net/http"
	"time"
)

type User struct {
	ID           string              // Unique User ID
	DisplayName  string              // User's display name
	MessageQueue chan Message        // Channel to receive messages
	SSEWriter    http.ResponseWriter // SSE connection for sending messages
}

type ChatRoom struct {
	ID        string           // Unique Room ID
	Name      string           // Room Name
	Members   map[string]*User // Active members (UserID -> User)
	Broadcast chan Message     // Broadcast message channel
}

type Message struct {
	SenderID   string    // User ID of the sender
	ReceiverID string    // Optional: For private messages
	RoomID     string    // Chat room ID (for broadcast messages)
	Content    string    // Message content
	Timestamp  time.Time // Time of the message
}
