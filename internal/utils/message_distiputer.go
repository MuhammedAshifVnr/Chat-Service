package utils

// import (
// 	"fmt"
// 	"log"

// 	"github.com/MuhammedAshifVnr/Chat-Service/internal/models"
// )

// StartRoomMessageDispatcher starts listening for broadcast messages in a room
// func StartRoomMessageDispatcher(room *models.ChatRoom) {
// 	log.Printf("Starting message dispatcher for room: %s", room.Name)

// 	// Infinite loop to listen for incoming messages
// 	for message := range room.Broadcast {
// 		log.Printf("Dispatching message in room %s: %v", room.Name, message)

// 		// Broadcast the message to all members of the room
// 		room.Members.Range(func(_, member interface{}) bool {
// 			fmt.Println("====")
// 			if user, ok := member.(*models.MemberInfo); ok {
// 				fmt.Println("===", user)
// 				select {
// 				case user.MessageQueue <- message: // Non-blocking send
// 					log.Printf("Message sent to user %s: %s", user.ID, message)
// 				default:
// 					log.Printf("User %s's message queue is full, dropping message", user.ID)
// 				}
// 			} else {
// 				log.Println("Encountered a non-user entity in the room members list")
// 			}
// 			return true // Continue iterating over Members
// 		})
// 	}

// 	// Room is being closed
// 	log.Printf("Stopping message dispatcher for room: %s", room.Name)
// }


