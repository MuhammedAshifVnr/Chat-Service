package main

import (
	"log"
	"net/http"

	"github.com/MuhammedAshifVnr/Chat-Service/internal/core"
	"github.com/MuhammedAshifVnr/Chat-Service/internal/handlers"
)

func main() {
	// Initialize core components
	roomManager := core.NewRoomManager()
	userManager := core.NewUserManager()
	messageDispatcher := core.NewMessageDispatcher(roomManager, userManager)

	// Initialize handlers
	chatRoomHandler := handlers.NewChatRoomHandler(roomManager)
	userHandler := handlers.NewUserHandler(userManager)
	messageHandler := handlers.NewMessageHandler(messageDispatcher, userManager, roomManager)

	// Set up routes
	mux := http.NewServeMux()

	// Chat room routes
	mux.HandleFunc("/rooms", chatRoomHandler.CreateRoomHandler)             // POST /rooms - Create a room
	mux.HandleFunc("/rooms/list", chatRoomHandler.ListRoomsHandler)        // GET /rooms/list - List all rooms
	mux.HandleFunc("/rooms/join", chatRoomHandler.JoinRoomHandler)         // POST /rooms/join - Join a room
	mux.HandleFunc("/rooms/leave", chatRoomHandler.LeaveRoomHandler)       // POST /rooms/leave - Leave a room
	mux.HandleFunc("/rooms/members", chatRoomHandler.ListMembersHandler)   // GET /rooms/members?room_id=<roomID> - List room members

	// User routes
	mux.HandleFunc("/users", userHandler.CreateUserHandler)                // POST /users - Create a user
	mux.HandleFunc("/users/get", userHandler.GetUserHandler)               // GET /users/get?id=<userID> - Get user details
	mux.HandleFunc("/users/update", userHandler.UpdateUserHandler)         // POST /users/update - Update user details
	mux.HandleFunc("/users/delete", userHandler.DeleteUserHandler)         // DELETE /users/delete?id=<userID> - Delete a user

	// Message routes
	mux.HandleFunc("/messages/broadcast", messageHandler.HandleBroadcastMessage) // POST /messages/broadcast - Broadcast message
	mux.HandleFunc("/messages/private", messageHandler.HandlePrivateMessage)     // POST /messages/private - Private message
	mux.HandleFunc("/sse", messageHandler.HandleSSEConnection)                   // GET /sse?user_id=<userID> - SSE connection

	// Start the server
	server := &http.Server{
		Addr:    ":8080", // Change to your preferred port
		Handler: mux,
	}

	log.Println("Chat service running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
