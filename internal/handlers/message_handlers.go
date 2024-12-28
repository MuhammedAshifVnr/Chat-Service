package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MuhammedAshifVnr/Chat-Service/internal/core"
	"github.com/MuhammedAshifVnr/Chat-Service/internal/utils"
)

type MessageHandler struct {
	MessageDispatcher *core.MessageDispatcher
	UserManager       *core.UserManager
	RoomManager       *core.RoomManager
}

func NewMessageHandler(md *core.MessageDispatcher, um *core.UserManager, rm *core.RoomManager) *MessageHandler {
	return &MessageHandler{
		MessageDispatcher: md,
		UserManager:       um,
		RoomManager:       rm,
	}
}

// HandleBroadcastMessage handles broadcasting a message to a chat room.
func (h *MessageHandler) HandleBroadcastMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID  string `json:"room_id"`
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoomID == "" || req.UserID == "" || req.Content == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.MessageDispatcher.BroadcastMessage(req.RoomID, req.UserID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandlePrivateMessage handles sending a private message between users.
func (h *MessageHandler) HandlePrivateMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SenderID   string `json:"sender_id"`
		ReceiverID string `json:"receiver_id"`
		Content    string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SenderID == "" || req.ReceiverID == "" || req.Content == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.MessageDispatcher.SendPrivateMessage(req.SenderID, req.ReceiverID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleSSEConnection handles the Server-Sent Events connection for real-time updates.
func (h *MessageHandler) HandleSSEConnection(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	// Fetch the user.
	user, err := h.UserManager.GetUser(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Configure headers for SSE.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Listen to the user's message queue.
	for {
		select {
		case msg, ok := <-user.MessageQueue:
			if !ok {
				// Channel closed, terminate SSE connection.
				return
			}
			// Write the message to the SSE stream.
			if err := utils.WriteSSE(w, "message", msg.Content); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
