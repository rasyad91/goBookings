package dbrepo

import (
	"fmt"
	"time"

	"github.com/rasyad91/goBookings/internal/models"
)

func (m *testDBrepo) GetAllUsers() bool {
	return true
}

// InsertReservation into the database
func (m *testDBrepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2 then fail; otherwise pass
	if res.RoomID == 2 {
		return 0, fmt.Errorf("some error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restrictions into the database
func (m *testDBrepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 100 {
		return fmt.Errorf("some error")
	}
	return nil
}

// IsAvailableByDatesByRoomID returns true if availability exists for roomID, and false if no availability exists
func (m *testDBrepo) IsAvailableByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms for the given dates
func (m *testDBrepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	return []models.Room{}, nil
}

// GetRoomByID gets room by ID
func (m *testDBrepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id < 1 {
		return room, fmt.Errorf("some error")
	}

	return models.Room{}, nil
}

// GetUserByID returns a user by ID
func (m *testDBrepo) GetUserByID(id int) (models.User, error) {

	return models.User{}, nil
}

// UpdateUser to the database
func (m *testDBrepo) UpdateUser(u models.User) error {

	return nil
}

func (m *testDBrepo) Authenticate(email, password string) (int, string, error) {

	return 0, "", nil
}
