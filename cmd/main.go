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
	chatRoomHandler := handlers.NewChatRoomHandler(roomManager, messageDispatcher, userManager)
	userHandler := handlers.NewUserHandler(userManager, roomManager)
	messageHandler := handlers.NewMessageHandler(messageDispatcher, userManager, roomManager)

	// Set up routes
	mux := http.NewServeMux()

	// Chat room routes
	mux.HandleFunc("/rooms", chatRoomHandler.CreateRoomHandler)          // POST /rooms - Create a room
	mux.HandleFunc("/rooms/list", chatRoomHandler.ListRoomsHandler)      // GET /rooms/list - List all rooms
	mux.HandleFunc("/rooms/join", chatRoomHandler.JoinRoomHandler)       // POST /rooms/join - Join a room
	mux.HandleFunc("/rooms/leave", chatRoomHandler.LeaveRoomHandler)     // POST /rooms/leave - Leave a room
	mux.HandleFunc("/rooms/members", chatRoomHandler.ListMembersHandler) // GET /rooms/members?room_id=<roomID> - List room members
	mux.HandleFunc("/rooms/delete", chatRoomHandler.DeleteRoomHandler)   //DELETE /rooms/delete -Delete a room

	// User routes
	mux.HandleFunc("/users", userHandler.CreateUserHandler)        // POST /users - Create a user
	mux.HandleFunc("/users/get", userHandler.GetUserHandler)       // GET /users/get?id=<userID> - Get user details
	mux.HandleFunc("/users/update", userHandler.UpdateUserHandler) // POST /users/update - Update user details
	mux.HandleFunc("/users/delete", userHandler.DeleteUserHandler) // DELETE /users/delete?id=<userID> - Delete a user
	mux.HandleFunc("/users/all", userHandler.GetAllUsersHandler)   // GET /users/all - Get all users

	// Message routes
	mux.HandleFunc("/messages/broadcast", messageHandler.HandleBroadcastMessage) // POST /messages/broadcast - Broadcast message
	mux.HandleFunc("/messages/private", messageHandler.HandlePrivateMessage)     // POST /messages/private - Private message
	mux.HandleFunc("/sse/broadcast", messageHandler.HandleSSEConnection)         // GET /sse/broadcast?user_id=<userID> - SSE connection for broadcast
	mux.HandleFunc("/sse/private", messageHandler.HandlePrivateSSEConnection)    //GET /sse/privet?user_id=<userID> -  SSE connection for privet

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
