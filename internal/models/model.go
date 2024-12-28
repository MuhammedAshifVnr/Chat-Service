package models

import (
	"net/http"
	"sync"
	"time"
)

type User struct {
	ID           string              // Unique User ID
	DisplayName  string              // User's display name
	MessageQueue chan Message        // Channel to receive messages
	SSEWriter    http.ResponseWriter // SSE connection for sending messages
}

type MemberInfo struct {
	UserID      string // Unique User ID
	DisplayName string // User's display name
}

type ChatRoom struct {
	ID        string       // Unique Room ID
	Name      string       // Room Name
	Members   sync.Map     // Thread-safe map of members (key: userID, value: MemberInfo)
	Broadcast chan Message // Broadcast message channel
}

type Message struct {
	SenderID   string    // User ID of the sender
	ReceiverID string    // Optional: For private messages
	RoomID     string    // Chat room ID (for broadcast messages)
	Content    string    // Message content
	Timestamp  time.Time // Time of the message
}

func (r *ChatRoom) ListMembers() []string {
	var members []string

	// Iterate through all members in the sync.Map.
	r.Members.Range(func(key, value interface{}) bool {
		members = append(members, value.(string)) // Assuming value is the display name.
		return true
	})

	return members
}
