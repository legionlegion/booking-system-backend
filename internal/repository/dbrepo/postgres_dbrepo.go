package dbrepo

import (
	"booking-backend/internal/models"
	"context"
	"database/sql"
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

func (m *PostgresDBRepo) AllBookings() ([]*models.Booking, error) {
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
		return nil, err
	}

	defer rows.Close()

	var bookings []*models.Booking

	for rows.Next() {
		var booking models.Booking
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

	return bookings, nil
}

func (m *PostgresDBRepo) InsertBookingRequest(booking models.Booking) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into requestedbookings (username, name, date, unit_number, start_time,
		end_time, purpose, facility)
		values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
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
		return 0, nil
	}

	return newID, nil
}

func (m *PostgresDBRepo) RegisterUser(username, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	stmt := `insert into users (username, password, is_admin) values ($1, $2, $3) returning id`

	var newID int

	err = m.DB.QueryRowContext(ctx, stmt, username, string(hashedPassword), false).Scan(&newID)

	if err != nil {
		return 0, err
	}

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
