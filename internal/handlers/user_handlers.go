package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MuhammedAshifVnr/Chat-Service/internal/core"
)

type UserHandler struct {
	UserManager *core.UserManager
}

func NewUserHandler(um *core.UserManager) *UserHandler {
	return &UserHandler{UserManager: um}
}

// CreateUserHandler handles user creation
func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user := uh.UserManager.AddUser(req.DisplayName)
	json.NewEncoder(w).Encode(user)
}

// GetUserHandler retrieves user details by ID
func (uh *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	user, err := uh.UserManager.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// UpdateUserHandler updates a user's display name
func (uh *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      string `json:"user_id"`
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := uh.UserManager.UpdateDisplayName(req.UserID, req.DisplayName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteUserHandler removes a user by ID
func (uh *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	err := uh.UserManager.RemoveUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
