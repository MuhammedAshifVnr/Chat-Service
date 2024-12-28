// core/message_dispatcher.go
package core

import (
	"errors"
	"sync"
	"time"

	"github.com/MuhammedAshifVnr/Chat-Service/internal/models"
)

type MessageDispatcher struct {
	RoomManager *RoomManager
	UserManager *UserManager
}

func NewMessageDispatcher(rm *RoomManager, um *UserManager) *MessageDispatcher {
	return &MessageDispatcher{
		RoomManager: rm,
		UserManager: um,
	}
}

// BroadcastMessage sends a message to all members of a room
func (md *MessageDispatcher) BroadcastMessage(roomID, senderID, content string) error {
	room, err := md.RoomManager.GetRoom(roomID)
	if err != nil {
		return err
	}

	message := models.Message{
		SenderID:  senderID,
		RoomID:    roomID,
		Content:   content,
		Timestamp: time.Now(),
	}

	room.Broadcast <- message
	return nil
}

// SendPrivateMessage sends a private message between two users
func (md *MessageDispatcher) SendPrivateMessage(senderID, receiverID, content string) error {
	sender, err := md.UserManager.GetUser(senderID)
	if err != nil {
		return errors.New("sender not found")
	}

	receiver, err := md.UserManager.GetUser(receiverID)
	if err != nil {
		return errors.New("receiver not found")
	}

	message := models.Message{
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
		Content:    content,
		Timestamp:  time.Now(),
	}

	select {
	case receiver.MessageQueue <- message:
		return nil
	default:
		return errors.New("receiver's message queue is full")
	}
}

// StartRoomMessageDispatcher starts listening for broadcast messages in a room
func (md *MessageDispatcher) StartRoomMessageDispatcher(roomID string, wg *sync.WaitGroup) {
	defer wg.Done()

	room, err := md.RoomManager.GetRoom(roomID)
	if err != nil {
		return
	}

	for message := range room.Broadcast {
		room.Members.Range(func(_, value interface{}) bool {
			member := value.(models.MemberInfo)
			user, err := md.UserManager.GetUser(member.UserID)
			if err == nil {
				select {
				case user.MessageQueue <- message:
				default:
					// Drop message if user's queue is full
				}
			}
			return true
		})
	}
}
