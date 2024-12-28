package core

import (
	"errors"
	"sync"

	"github.com/MuhammedAshifVnr/Chat-Service/internal/models"
	"github.com/google/uuid"
)

type UserManager struct {
	Users sync.Map // Thread-safe map to store users
}

func NewUserManager() *UserManager {
	return &UserManager{}
}

// AddUser adds a new user and returns the user object
func (um *UserManager) AddUser(displayName string) *models.User {
	userID := uuid.New().String()
	user := &models.User{
		ID:           userID,
		DisplayName:  displayName,
		MessageQueue: make(chan models.Message, 1000),
	}
	um.Users.Store(userID, user)
	return user
}

// GetUser fetches a user by ID
func (um *UserManager) GetUser(userID string) (*models.User, error) {
	user, ok := um.Users.Load(userID)
	if !ok {
		return nil, errors.New("user not found")
	}
	return user.(*models.User), nil
}

// RemoveUser removes a user by ID and closes their message queue
func (um *UserManager) RemoveUser(userID string) error {
	user, ok := um.Users.LoadAndDelete(userID)
	if !ok {
		return errors.New("user not found")
	}
	close(user.(*models.User).MessageQueue)
	return nil
}

// UpdateDisplayName updates a user’s display name
func (um *UserManager) UpdateDisplayName(userID string, newName string) error {
	user, err := um.GetUser(userID)
	if err != nil {
		return err
	}
	user.DisplayName = newName
	return nil
}

// DisconnectUser handles user cleanup when they disconnect
func (um *UserManager) DisconnectUser(userID string) error {
	return um.RemoveUser(userID)
}
