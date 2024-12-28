package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/MuhammedAshifVnr/Chat-Service/internal/core"
	"github.com/MuhammedAshifVnr/Chat-Service/internal/handlers"
)

func main() {
	// Initialize managers
	userManager := core.NewUserManager()
	roomManager := core.NewRoomManager()
	messageDispatcher := core.NewMessageDispatcher(roomManager, userManager)

	// Start dispatcher for each room
	go startRoomDispatchers(roomManager, messageDispatcher)

	// Initialize router
	router := mux.NewRouter()

	// Register handlers
	handlers.RegisterUserHandlers(router, userManager)
	handlers.RegisterRoomHandlers(router, roomManager)
	handlers.RegisterMessageHandlers(router, messageDispatcher)

	// Start HTTP server
	log.Println("Chat server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Helper to start dispatchers for each room
func startRoomDispatchers(roomManager *core.RoomManager, dispatcher *core.MessageDispatcher) {
	var wg sync.WaitGroup
	roomManager.Rooms.Range(func(_, value interface{}) bool {
		room := value.(*models.ChatRoom)
		wg.Add(1)
		go dispatcher.StartRoomMessageDispatcher(room.ID, &wg)
		return true
	})
	wg.Wait()
}
