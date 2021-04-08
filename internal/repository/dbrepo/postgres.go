package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/rasyad91/goBookings/internal/models"
)

func (m *postgresDBrepo) GetAllUsers() bool {
	return true
}

// InsertReservation into the database
func (m *postgresDBrepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `insert into reservations 
			(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restrictions into the database
func (m *postgresDBrepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions 
			(start_date,end_date, room_id, reservation_id, created_at, updated_at, restriction_id)
			values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}

	return nil
}

// IsAvailableByDatesByRoomID returns true if availability exists for roomID, and false if no availability exists
func (m *postgresDBrepo) IsAvailableByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			select 
				count(id)
			from 
				room_restrictions
			where 
				$1 < end_date and $2 > start_date
				and room_id = $3`

	var numRows int

	row := m.DB.QueryRowContext(ctx, query, end, start, roomID)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	return numRows == 0, nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms for the given dates
func (m *postgresDBrepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println(start, end)

	// arbitary number 10 to make rooms
	rooms := make([]models.Room, 0, 10)

	query := `
			select 
				r.id , r.room_name 
			from 
				rooms r
			where 
				r.id not in 
				(select room_id from room_restrictions where $1 < end_date and $2 > start_date)
			`

	rows, err := m.DB.QueryContext(ctx, query, end, start)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		rows.Scan(
			&room.ID,
			&room.RoomName)

		if err := rows.Err(); err != nil {
			return rooms, fmt.Errorf("error querying DB: SearchAvailabilityForAllRooms: %w", err)
		}

		rooms = append(rooms, room)
	}
	return rooms, nil
}

// GetRoomByID gets room by ID
func (m *postgresDBrepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, room_name from rooms where id = $1`

	room := models.Room{}

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName)
	if err != nil {
		return room, err
	}

	return room, nil
}