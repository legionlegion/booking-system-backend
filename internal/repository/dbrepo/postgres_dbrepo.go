package dbrepo

import (
	"booking-backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) AllBookings() ([]*models.ApprovedBooking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			id, username, name, date, unit_number, start_time, end_time, purpose, facility 
		from 
			approvedbookings 
		order by 
			id
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Print("Err in AllBookings: ", err)
		return nil, err
	}

	defer rows.Close()

	var bookings []*models.ApprovedBooking
	for rows.Next() {
		var booking models.ApprovedBooking
		err := rows.Scan(
			&booking.ID,
			&booking.Username,
			&booking.Name,
			&booking.Date,
			&booking.UnitNumber,
			&booking.StartTime,
			&booking.EndTime,
			&booking.Purpose,
			&booking.Facility,
		)

		if err != nil {
			return nil, err
		}

		bookings = append(bookings, &booking)
	}
	log.Print("Bookings: ", bookings)
	return bookings, nil
}

func (m *PostgresDBRepo) InsertBookingRequest(booking models.Booking) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	log.Print("Booking: ", booking)
	defer cancel()

	// check for overlaps, WHERE clause covers all overlap scenarios
	checkOverlapStmt := `
	SELECT id
	FROM approvedbookings
	WHERE 
		($1 < end_time AND $2 > start_time)
	LIMIT 1;
	`

	var overlapID int
	err := m.DB.QueryRowContext(ctx, checkOverlapStmt, booking.StartTime, booking.EndTime).Scan(&overlapID)
	if err != nil && err != sql.ErrNoRows {
		log.Print("Overlap check error: ", err)
		return 0, err
	}
	log.Print("Overlap ID: ", overlapID)

	// If overlap found, return error
	if err != sql.ErrNoRows {
		return 0, errors.New("Booking time overlaps with an existing booking")
	}

	// If no overlaps, proceed with insertion
	stmt := `insert into approvedbookings (username, name, date, unit_number, start_time,
		end_time, purpose, facility)
		values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var newID int

	err = m.DB.QueryRowContext(ctx, stmt,
		booking.Username,
		booking.Name,
		booking.Date,
		booking.UnitNumber,
		booking.StartTime,
		booking.EndTime,
		booking.Purpose,
		booking.Facility,
	).Scan(&newID)

	if err != nil {
		log.Print("Insertion err: ", err)
		return 0, err
	}
	log.Print("New id: ", newID)

	return newID, nil
}

func (m *PostgresDBRepo) ApproveBooking(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Copy the booking from requestedbookings to approvedbookings
	copyStmt := `INSERT INTO approvedbookings (username, name, date, unit_number, start_time, end_time, purpose, facility)
		SELECT username, name, date, unit_number, start_time, end_time, purpose, facility FROM requestedbookings WHERE id = $1`
	_, err = tx.ExecContext(ctx, copyStmt, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Delete the booking from requestedbookings
	deleteStmt := `DELETE FROM requestedbookings WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteStmt, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetUserByName(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, username, password from users where username = $1`
	var user models.User
	row := m.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}


func (m *PostgresDBRepo) RegisterUser(username, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	stmt := `insert into users (username, password, is_admin) values ($1, $2, $3) returning id`

	var newID int

	err = m.DB.QueryRowContext(ctx, stmt, username, string(hashedPassword), false).Scan(&newID)

	if err != nil {
		return nil, err
	}

	var user models.User = models.User{
		ID: newID,
		Username: username,
		Password: password,
	}

	return &user, nil
}