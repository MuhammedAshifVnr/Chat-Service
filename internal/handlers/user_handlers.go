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
	RoomManager *core.RoomManager
}

// NewUserHandler initializes a new UserHandler.
func NewUserHandler(um *core.UserManager, rm *core.RoomManager) *UserHandler {
	return &UserHandler{UserManager: um, RoomManager: rm}
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

	user, err := uh.UserManager.AddUser(req.DisplayName)
	if err != nil {
		log.Printf("display name already exist: %v", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "display name already taken"})
		return
	}
	log.Printf("User created with ID: %s, DisplayName: %s", user.ID, user.DisplayName)
	response := struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	}{
		ID:          user.ID,
		DisplayName: user.DisplayName,
	}
	respondJSON(w, http.StatusCreated, response)
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
	response := struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	}{
		ID:          user.ID,
		DisplayName: user.DisplayName,
	}
	log.Printf("User details fetched for ID: %s", userID)
	respondJSON(w, http.StatusOK, response)
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
	user, err := uh.UserManager.GetUser(userID)
	if err != nil {
		log.Printf("User not found: %s", userID)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}
	if user.RoomIn != "" {
		room, _ := uh.RoomManager.GetRoom(user.RoomIn)
		room.RemoveMember(userID)
	}
	err = uh.UserManager.RemoveUser(userID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	log.Printf("User deleted: ID %s", userID)
	respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (uh *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Log the request for debugging purposes
	log.Println("Received request to fetch all users")
	// Fetch all users from UserManager
	users, err := uh.UserManager.GetAllUsers()
	if err != nil {
		log.Println("No users found")
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "No users found"})
		return
	}

	// Transform the user data into a response format
	response := make([]map[string]string, 0, len(users))
	for _, user := range users {
		response = append(response, map[string]string{
			"id":           user.ID,
			"display_name": user.DisplayName,
		})
	}

	// Log successful retrieval
	log.Printf("Fetched %d users", len(response))

	// Send response with status 200
	respondJSON(w, http.StatusOK, response)
}
