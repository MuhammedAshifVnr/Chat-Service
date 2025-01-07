package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MuhammedAshifVnr/Chat-Service/internal/core"
	"github.com/MuhammedAshifVnr/Chat-Service/internal/utils"
)

// MessageHandler handles messaging-related operations.
type MessageHandler struct {
	MessageDispatcher *core.MessageDispatcher
	UserManager       *core.UserManager
	RoomManager       *core.RoomManager
}

// NewMessageHandler initializes a new MessageHandler.
func NewMessageHandler(md *core.MessageDispatcher, um *core.UserManager, rm *core.RoomManager) *MessageHandler {
	return &MessageHandler{
		MessageDispatcher: md,
		UserManager:       um,
		RoomManager:       rm,
	}
}

// HandleBroadcastMessage handles broadcasting a message to a chat room.
func (h *MessageHandler) HandleBroadcastMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to broadcast a message")

	var req struct {
		RoomID  string `json:"room_id"`
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoomID == "" || req.UserID == "" || req.Content == "" {
		log.Printf("Invalid broadcast message request: %v", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err := h.MessageDispatcher.BroadcastMessage(req.RoomID, req.UserID, req.Content)
	if err != nil {
		log.Printf("Failed to broadcast message: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Message broadcasted to room %s by user %s", req.RoomID, req.UserID)
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Message broadcasted successfully",
	})
}

// HandlePrivateMessage handles sending a private message between users.
func (h *MessageHandler) HandlePrivateMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to send a private message")

	var req struct {
		SenderID   string `json:"sender_id"`
		ReceiverID string `json:"receiver_id"`
		Content    string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SenderID == "" || req.ReceiverID == "" || req.Content == "" {
		log.Printf("Invalid private message request: %v", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err := h.MessageDispatcher.SendPrivateMessage(req.SenderID, req.ReceiverID, req.Content)
	if err != nil {
		log.Printf("Failed to send private message: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Private message sent from user %s to user %s", req.SenderID, req.ReceiverID)
	respondJSON(w, http.StatusOK, map[string]string{"message": "Private message sent successfully"})
}

// HandleSSEConnection handles the Server-Sent Events connection for real-time updates.
func (h *MessageHandler) HandleSSEConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to establish an SSE connection")

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Println("Missing user_id parameter in SSE connection")
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing user_id parameter"})
		return
	}

	// Fetch the user.
	user, err := h.UserManager.GetUser(userID)
	if err != nil {
		log.Printf("User not found for SSE connection: %s", userID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	// Configure headers for SSE.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("Streaming not supported in SSE connection")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	log.Printf("SSE connection established for user %s", userID)

	// Listen to the user's message queue.
	for {
		select {
		case msg, ok := <-user.MessageQueue:
			if !ok {
				log.Printf("Message queue closed for user %s", userID)
				return
			}
			// Write the message to the SSE stream.
			if err := utils.WriteSSE(w, msg.SenderID, msg.Content, msg.Timestamp.String()); err != nil {
				log.Printf("Error writing SSE message for user %s: %v", userID, err)
				return
			}
			flusher.Flush()
		}
	}
}

func (h *MessageHandler) HandlePrivateSSEConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to establish a private SSE connection")

	// Get the user_id from the URL parameters
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Println("Missing user_id parameter in SSE connection")
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing user_id parameter"})
		return
	}

	// Fetch the user.
	user, err := h.UserManager.GetUser(userID)
	if err != nil {
		log.Printf("User not found for SSE connection: %s", userID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	// Configure headers for SSE.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("Streaming not supported in SSE connection")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Listen to the user's private message queue.
	log.Printf("Private SSE connection established for user %s", userID)

	for {
		select {
		case msg, ok := <-user.PrivateMessageQueue:
			if !ok {
				log.Printf("Message queue closed for user %s", userID)
				return
			}
			// Write the message to the SSE stream
			if err := utils.WriteSSE(w, msg.SenderID, msg.Content, msg.Timestamp.String()); err != nil {
				log.Printf("Error writing SSE message for user %s: %v", userID, err)
				return
			}
			flusher.Flush()
		}
	}
}
