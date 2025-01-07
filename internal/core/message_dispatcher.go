package core

import (
	"fmt"
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
	sender, err := md.UserManager.GetUser(senderID)
	if err != nil {
		return fmt.Errorf("receiver not found: %v", err)
	}
	message := models.Message{
		SenderID:  sender.DisplayName,
		RoomID:    roomID,
		Content:   content,
		Timestamp: time.Now(),
	}

	room.Broadcast <- message
	return nil
}

// SendPrivateMessage sends a private message between two users

func (md *MessageDispatcher) SendPrivateMessage(senderID, receiverID, content string) error {
	_, err := md.UserManager.GetUser(senderID)
	if err != nil {
		return fmt.Errorf("sender not found: %v", err)
	}

	receiver, err := md.UserManager.GetUser(receiverID)
	if err != nil {
		return fmt.Errorf("receiver not found: %v", err)
	}
	sender, err := md.UserManager.GetUser(senderID)
	if err != nil {
		return fmt.Errorf("receiver not found: %v", err)
	}
	message := models.Message{
		SenderID: sender.DisplayName,
		Content:  content,
	}

	select {
	case receiver.PrivateMessageQueue <- message:
		return nil
	default:
		return fmt.Errorf("receiver's private message queue is full")
	}
}

// StartRoomMessageDispatcher starts listening for broadcast messages in a room
func (md *MessageDispatcher) StartRoomMessageDispatcher(roomID string) {
	room, err := md.RoomManager.GetRoom(roomID)
	if err != nil {
		return
	}

	numWorkers := 5
	workerDone := make(chan struct{}, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go md.startRoomWorker(room, workerDone)
	}

	for i := 0; i < numWorkers; i++ {
		<-workerDone
	}
}

// Worker function that listens to the room's broadcast channel
func (md *MessageDispatcher) startRoomWorker(room *ChatRoom, done chan struct{}) {
	defer func() { done <- struct{}{} }() // Signal that the worker is done

	for {
		select {
		case message, ok := <-room.Broadcast:
			if !ok {
				return // Channel closed, stop the worker
			}

			// Distribute the message to each member of the room
			room.Members.Range(func(_, value interface{}) bool {
				member := value.(models.MemberInfo)
				user, err := md.UserManager.GetUser(member.UserID)
				if err == nil {
					select {
					case user.MessageQueue <- message:
					default:
						// Drop message if the user's queue is full
					}
				}
				return true
			})

		case <-room.Done:
			return // Stop the worker when the room is deleted
		}
	}
}

// func (md *MessageDispatcher) StartRoomMessageDispatcher(roomID string) {
// 	room, err := md.RoomManager.GetRoom(roomID)
// 	if err != nil {
// 		return
// 	}

// 	for {
// 		select {
// 		case message, ok := <-room.Broadcast:
// 			if !ok {
// 				return // Channel closed, stop the dispatcher
// 			}
// 			room.Members.Range(func(_, value interface{}) bool {
// 				member := value.(models.MemberInfo)
// 				user, err := md.UserManager.GetUser(member.UserID)
// 				if err == nil {
// 					select {
// 					case user.MessageQueue <- message:
// 					default:
// 						// Drop message if user's queue is full
// 					}
// 				}
// 				return true
// 			})
// 		case <-room.Done:
// 			return // Stop when the room is deleted
// 		}
// 	}
// }
