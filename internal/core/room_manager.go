package core

import (
	"errors"
	"sync"
)

type RoomManager struct {
	Rooms sync.Map // Thread-safe map to store rooms
}

func NewRoomManager() *RoomManager {
	return &RoomManager{}
}

// CreateRoom creates a new chat room with the given name.
func (rm *RoomManager) CreateRoom(name string) (*ChatRoom, error) {
	if name == "" {
		return nil, errors.New("room name cannot be empty")
	}

	newRoom := NewChatRoom(name, name)

	_, loaded := rm.Rooms.LoadOrStore(name, newRoom)
	if loaded {
		return nil, errors.New("room already exists")
	}
	return newRoom, nil
}

// GetRoom fetches a chat room by name.
func (rm *RoomManager) GetRoom(name string) (*ChatRoom, error) {
	if room, ok := rm.Rooms.Load(name); ok {
		return room.(*ChatRoom), nil
	}
	return nil, errors.New("room not found")
}

// ListRooms lists all available chat rooms.
func (rm *RoomManager) ListRooms() []string {
	var roomNames []string
	rm.Rooms.Range(func(key, _ interface{}) bool {
		roomNames = append(roomNames, key.(string))
		return true
	})
	return roomNames
}

// DeleteRoom deletes a room by name if it exists and is empty.
func (rm *RoomManager) DeleteRoom(name string) error {
	room, err := rm.GetRoom(name)
	if err != nil {
		return err
	}

	var isEmpty = true
	room.Members.Range(func(_, _ interface{}) bool {
		isEmpty = false
		return false
	})

	if !isEmpty {
		return errors.New("cannot delete room with active members")
	}

	rm.Rooms.Delete(name)
	close(room.Broadcast)
	return nil
}
