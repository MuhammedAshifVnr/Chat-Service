package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MuhammedAshifVnr/Chat-Service/internal/core"
)

// ChatRoomHandler holds the RoomManager instance for managing chat rooms.
type ChatRoomHandler struct {
	RoomManager *core.RoomManager
}

// NewChatRoomHandler initializes a new ChatRoomHandler.
func NewChatRoomHandler(roomManager *core.RoomManager) *ChatRoomHandler {
	return &ChatRoomHandler{
		RoomManager: roomManager,
	}
}

// respondJSON sends a JSON response with the provided status code and data.
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// CreateRoomHandler handles the creation of a new chat room.
func (h *ChatRoomHandler) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a new room")

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		log.Printf("Invalid room name: %v", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid room name"})
		return
	}

	room, err := h.RoomManager.CreateRoom(req.Name)
	if err != nil {
		log.Printf("Failed to create room: %v", err)
		respondJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Room created successfully: %s", room.ID)
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"room_id":  room.ID,
		"name":     room.Name,
		"message":  "Room created successfully",
	})
}

// ListRoomsHandler lists all existing chat rooms.
func (h *ChatRoomHandler) ListRoomsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to list all rooms")

	rooms := h.RoomManager.ListRooms()
	log.Printf("Rooms found: %d", len(rooms))
	respondJSON(w, http.StatusOK, rooms)
}

// JoinRoomHandler allows a user to join a chat room.
func (h *ChatRoomHandler) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to join a room")

	var req struct {
		RoomID      string `json:"room_id"`
		UserID      string `json:"user_id"`
		DisplayName string `json:"display_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoomID == "" || req.UserID == "" || req.DisplayName == "" {
		log.Printf("Invalid input for joining room: %v", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		return
	}

	room, err := h.RoomManager.GetRoom(req.RoomID)
	if err != nil {
		log.Printf("Room not found: %s", req.RoomID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Room not found"})
		return
	}

	room.AddMember(req.UserID, req.DisplayName)
	log.Printf("User %s joined room %s", req.DisplayName, req.RoomID)
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "User joined the room successfully",
	})
}

// LeaveRoomHandler allows a user to leave a chat room.
func (h *ChatRoomHandler) LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to leave a room")

	var req struct {
		RoomID string `json:"room_id"`
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoomID == "" || req.UserID == "" {
		log.Printf("Invalid input for leaving room: %v", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		return
	}

	room, err := h.RoomManager.GetRoom(req.RoomID)
	if err != nil {
		log.Printf("Room not found: %s", req.RoomID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Room not found"})
		return
	}

	room.RemoveMember(req.UserID)
	log.Printf("User %s left room %s", req.UserID, req.RoomID)
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "User left the room successfully",
	})
}

// ListMembersHandler lists all members in a chat room.
func (h *ChatRoomHandler) ListMembersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to list members in a room")

	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		log.Println("Room ID not provided")
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Room ID is required"})
		return
	}

	room, err := h.RoomManager.GetRoom(roomID)
	if err != nil {
		log.Printf("Room not found: %s", roomID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Room not found"})
		return
	}

	members := room.ListMembers()
	log.Printf("Members in room %s: %d", roomID, len(members))
	respondJSON(w, http.StatusOK, members)
}
