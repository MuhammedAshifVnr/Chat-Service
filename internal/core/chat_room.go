package core

import (
	"sync"

	"github.com/MuhammedAshifVnr/Chat-Service/internal/models"
)

type ChatRoom struct {
	models.ChatRoom // Embedding the ChatRoom model
}

// NewChatRoom creates a new chat room instance
func NewChatRoom(id, name string) *ChatRoom {
	return &ChatRoom{
		ChatRoom: models.ChatRoom{
			ID:        id,
			Name:      name,
			Members:   sync.Map{},
			Broadcast: make(chan models.Message, 100), // Buffered channel for efficient broadcasting
		},
	}
}

// AddMember adds a user to the chat room
func (cr *ChatRoom) AddMember(userID, displayName string) {
	cr.Members.Store(userID, models.MemberInfo{
		UserID:      userID,
		DisplayName: displayName,
	})
}

// RemoveMember removes a user from the chat room
func (cr *ChatRoom) RemoveMember(userID string) {
	cr.Members.Delete(userID)
}

// ListMembers returns a list of all members in the chat room
func (cr *ChatRoom) ListMembers() []models.MemberInfo {
	members := []models.MemberInfo{}
	cr.Members.Range(func(_, value interface{}) bool {
		members = append(members, value.(models.MemberInfo))
		return true
	})
	return members
}

// SendBroadcast sends a message to the chat room's broadcast channel
func (cr *ChatRoom) SendBroadcast(message models.Message) {
	select {
	case cr.Broadcast <- message:
		// Message sent successfully
	default:
		// Drop the message if the broadcast channel is full
	}
}
