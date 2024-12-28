package handlers

import (
	"encoding/json"
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

// CreateRoomHandler handles the creation of a new chat room.
func (h *ChatRoomHandler) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	// Parse the request body.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid room name", http.StatusBadRequest)
		return
	}

	room, err := h.RoomManager.CreateRoom(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// Respond with the created room details.
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

// ListRoomsHandler lists all existing chat rooms.
func (h *ChatRoomHandler) ListRoomsHandler(w http.ResponseWriter, r *http.Request) {
	rooms := h.RoomManager.ListRooms()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

// JoinRoomHandler allows a user to join a chat room.
func (h *ChatRoomHandler) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID      string `json:"room_id"`
		UserID      string `json:"user_id"`
		DisplayName string `json:"display_name"`
	}

	// Parse the request body.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoomID == "" || req.UserID == "" || req.DisplayName == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	room, err := h.RoomManager.GetRoom(req.RoomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Add the user to the room.
	room.AddMember(req.UserID, req.DisplayName)

	// Respond with success.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User joined the room successfully",
	})
}

// LeaveRoomHandler allows a user to leave a chat room.
func (h *ChatRoomHandler) LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID string `json:"room_id"`
		UserID string `json:"user_id"`
	}

	// Parse the request body.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoomID == "" || req.UserID == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	room, err := h.RoomManager.GetRoom(req.RoomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Remove the user from the room.
	room.RemoveMember(req.UserID)

	// Respond with success.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User left the room successfully",
	})
}

// ListMembersHandler lists all members in a chat room.
func (h *ChatRoomHandler) ListMembersHandler(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	room, err := h.RoomManager.GetRoom(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	members := room.ListMembers()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(members)
}
