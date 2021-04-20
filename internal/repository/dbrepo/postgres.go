package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/rasyad91/goBookings/internal/models"
	"golang.org/x/crypto/bcrypt"
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
		return 0, fmt.Errorf("InsertReservation: %w", err)

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
		return fmt.Errorf("InsertRoomRestriction: %w", err)

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
		return false, fmt.Errorf("IsAvailableByDatesByRoomID: %w", err)

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

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, fmt.Errorf("SearchAvailabilityForAllRooms: %w", err)

	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		rows.Scan(
			&room.ID,
			&room.RoomName)

		if err := rows.Err(); err != nil {
			return rooms, fmt.Errorf("SearchAvailabilityForAllRooms: %w", err)
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
		return room, fmt.Errorf("GetRoomByID: %w", err)

	}

	return room, nil
}

// GetUserByID returns a user by ID
func (m *postgresDBrepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, email, first_name, password, access_level from users where id = $1`

	user := models.User{}

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.AccessLevel,
	)
	if err != nil {
		return user, fmt.Errorf("GetUserByID: %w", err)

	}

	return user, err
}

// UpdateUser to the database
func (m *postgresDBrepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users 
			set 
				first_name = $1, last_name=$2, email=$3, access_level = $4, updated_at = $5
			where 
			id = $5`
	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateUser: %w", err)

	}

	return nil
}

func (m *postgresDBrepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", fmt.Errorf("Authenticate: %w", err)

	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", fmt.Errorf("incorrect Password")
	}
	if err != nil {
		return 0, "", fmt.Errorf("Authenticate: %w", err)
	}

	return id, hashedPassword, nil
}

// GetAllReservations returns a slice reservations from the db
func (m *postgresDBrepo) GetAllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reservations := []models.Reservation{}

	query := ` 
			select
				r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id,
				r.created_at, r.updated_at, rm.room_name, r.processed
			from
				reservations r
			left join rooms rm on(r.room_id = rm.id)
			order by r.start_date desc
		`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return []models.Reservation{}, fmt.Errorf("GetAllReservations: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomID,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Room.RoomName,
			&reservation.Processed,
		)
		if err != nil {
			return reservations, fmt.Errorf("GetAllReservations: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return []models.Reservation{}, fmt.Errorf("GetAllReservations: %w", err)
	}

	return reservations, nil
}

// GetAllReservations returns a slice reservations from the db
func (m *postgresDBrepo) GetNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reservations := []models.Reservation{}

	query := ` 
			select
				r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id,
				r.created_at, r.updated_at, rm.room_name
			from
				reservations r
			left join rooms rm on(r.room_id = rm.id)
			where
				processed = false
			order by r.start_date desc
		`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return []models.Reservation{}, fmt.Errorf("GetNewReservations: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomID,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Room.RoomName,
		)
		if err != nil {
			return reservations, fmt.Errorf("GetNewReservations: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return []models.Reservation{}, fmt.Errorf("GetNewReservations: %w", err)
	}

	return reservations, nil
}

// GetReservationByID returns 1 reservation by ID
func (m *postgresDBrepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r := models.Reservation{}

	query := ` 
			select
				r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id,
				r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
			from
				reservations r
			left join rooms rm on (r.room_id = rm.id)
			where
				r.id = $1 
		`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&r.ID,
		&r.FirstName,
		&r.LastName,
		&r.Email,
		&r.Phone,
		&r.StartDate,
		&r.EndDate,
		&r.RoomID,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.Processed,
		&r.Room.ID,
		&r.Room.RoomName,
	)

	if err != nil {
		return r, fmt.Errorf("GetReservationByID: %w", err)
	}
	return r, nil
}

func (m *postgresDBrepo) UpdateReservation(r models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			update reservations set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
			where id = $6		
	`
	_, err := m.DB.ExecContext(ctx, query,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		time.Now(),
		r.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateReservation: %w", err)
	}

	return nil
}

func (m *postgresDBrepo) DeleteReservationByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			delete from reservations where id = $1		
			`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeleteReservationByID: %w", err)
	}
	return nil

}

func (m *postgresDBrepo) ProcessReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			update reservations set processed = true
			where id = $1		
	`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ProcessReservation: %w", err)
	}

	return nil
}

// GetAllReservations returns a slice reservations from the db
func (m *postgresDBrepo) GetAllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rooms := []models.Room{}

	query := ` 
			select
				r.id, r.room_name, r.created_at, r.updated_at
			from
				rooms r
		`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return []models.Room{}, fmt.Errorf("GetAllRooms: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var r models.Room
		err := rows.Scan(
			&r.ID,
			&r.RoomName,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return rooms, fmt.Errorf("GetAllRooms: %w", err)
		}
		rooms = append(rooms, r)
	}

	if err := rows.Err(); err != nil {
		return []models.Room{}, fmt.Errorf("GetAllRooms: %w", err)
	}

	return rooms, nil
}

// GetRestrictionsForRoomByDate returns restrictions for rooms by date range
func (m *postgresDBrepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	restrictions := []models.RoomRestriction{}

	query := `
			select 
				id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date
			from 
				room_restrictions
			where
			$1 < end_date and $2 >= start_date
			and room_id = $3
			
	`
	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return restrictions, fmt.Errorf("GetRestrictionsForRoomByDate: %w", err)
	}
	for rows.Next() {
		r := models.RoomRestriction{}
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return restrictions, fmt.Errorf("GetRestrictionsForRoomByDate: %w", err)
		}

		restrictions = append(restrictions, r)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return restrictions, fmt.Errorf("GetRestrictionsForRoomByDate: %w", err)
	}
	return restrictions, nil
}

func (m *postgresDBrepo) InsertBlockForRoom(roomID int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println(roomID, startDate)
	query := `
			insert into room_restrictions (room_id, start_date, end_date, restriction_id, created_at, updated_at)
			values ($1,$2,$3,$4,$5,$6)
			`
	_, err := m.DB.ExecContext(ctx, query,
		roomID,
		startDate,
		startDate.AddDate(0, 0, 1),
		2,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("InsertBlockForRoom: %w", err)

	}
	return nil
}

func (m *postgresDBrepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			delete from room_restrictions where id = $1
			`
	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("DeleteBlockByID: %w", err)

	}
	return nil
}
