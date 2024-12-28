package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MuhammedAshifVnr/Chat-Service/internal/core"
)

// UserHandler manages user-related operations.
type UserHandler struct {
	UserManager *core.UserManager
}

// NewUserHandler initializes a new UserHandler.
func NewUserHandler(um *core.UserManager) *UserHandler {
	return &UserHandler{UserManager: um}
}

// CreateUserHandler handles user creation.
func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a new user")

	var req struct {
		DisplayName string `json:"display_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DisplayName == "" {
		log.Printf("Invalid user creation request: %v", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body or missing display name"})
		return
	}

	user := uh.UserManager.AddUser(req.DisplayName)
	log.Printf("User created with ID: %s, DisplayName: %s", user.ID, user.DisplayName)

	respondJSON(w, http.StatusCreated, user)
}

// GetUserHandler retrieves user details by ID.
func (uh *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to fetch user details")

	userID := r.URL.Query().Get("id")
	if userID == "" {
		log.Println("Missing user ID in request")
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing user ID"})
		return
	}

	user, err := uh.UserManager.GetUser(userID)
	if err != nil {
		log.Printf("User not found: %s", userID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	log.Printf("User details fetched for ID: %s", userID)
	respondJSON(w, http.StatusOK, user)
}

// UpdateUserHandler updates a user's display name.
func (uh *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to update user display name")

	var req struct {
		UserID      string `json:"user_id"`
		DisplayName string `json:"display_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.DisplayName == "" {
		log.Printf("Invalid user update request: %v", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body or missing fields"})
		return
	}

	err := uh.UserManager.UpdateDisplayName(req.UserID, req.DisplayName)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	log.Printf("User updated: ID %s, New DisplayName: %s", req.UserID, req.DisplayName)
	respondJSON(w, http.StatusOK, map[string]string{"message": "User display name updated successfully"})
}

// DeleteUserHandler removes a user by ID.
func (uh *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to delete a user")

	userID := r.URL.Query().Get("id")
	if userID == "" {
		log.Println("Missing user ID in delete request")
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing user ID"})
		return
	}

	err := uh.UserManager.RemoveUser(userID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	log.Printf("User deleted: ID %s", userID)
	respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
