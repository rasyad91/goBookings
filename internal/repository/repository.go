package repository

import (
	"time"

	"github.com/rasyad91/goBookings/internal/models"
)

type DatabaseRepo interface {
	GetAllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	IsAvailableByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, password string) (int, string, error)
	GetAllReservations() ([]models.Reservation, error)
	GetNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(r models.Reservation) error
	DeleteReservationByID(id int) error
	ProcessReservation(id int) error
	GetAllRooms() ([]models.Room, error)
	GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error)
	InsertBlockForRoom(roomID int, startDate time.Time) error
	DeleteBlockByID(id int) error
}
