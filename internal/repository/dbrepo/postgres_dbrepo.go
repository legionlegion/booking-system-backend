package dbrepo

import (
	"booking-backend/internal/models"
	"context"
	"database/sql"
	"errors"
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

func (m *PostgresDBRepo) AllBookings() ([]*models.SubmittedBooking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select 
		id, username, name, start_date, end_date, unit_number, 
		start_time, end_time, purpose, facility 
	from 
		approvedbookings 
	where 
		start_time >= date_trunc('week', current_date) 
		and end_time < date_trunc('week', current_date) + interval '2 weeks'
	order by 
		id
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var bookings []*models.SubmittedBooking
	for rows.Next() {
		var booking models.SubmittedBooking
		err := rows.Scan(
			&booking.ID,
			&booking.Username,
			&booking.Name,
			&booking.StartDate,
			&booking.EndDate,
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

func (m *PostgresDBRepo) TwoWeekBookings() ([]*models.SubmittedBooking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select 
		id, username, name, start_date, end_date, unit_number, 
		start_time, end_time, purpose, facility 
	from 
		approvedbookings 
	where 
		start_time >= date_trunc('week', current_date) 
		and end_time < date_trunc('week', current_date) + interval '2 weeks'
	order by 
		id
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var bookings []*models.SubmittedBooking
	for rows.Next() {
		var booking models.SubmittedBooking
		err := rows.Scan(
			&booking.ID,
			&booking.Username,
			&booking.Name,
			&booking.StartDate,
			&booking.EndDate,
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

func (m *PostgresDBRepo) ManageBookings(username string) ([]*models.SubmittedBooking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	user, err := m.GetUserByName(username)
	if err != nil {
		return nil, err
	}

	var query string
	var rows *sql.Rows
	if user.IsAdmin {
		// get all bookings
		query = `
		select 
			id, username, name, start_date, end_date, unit_number, 
			start_time, end_time, purpose, facility 
		from 
			requestedbookings 
		order by 
			id
		`
		rows, err = m.DB.QueryContext(ctx, query)
	} else {
		// get bookings that belong to user only
		query = `
		select 
			id, username, name, start_date, end_date, unit_number, 
			start_time, end_time, purpose, facility 
		from 
			requestedbookings
		where
			username = $1
		order by 
			id
		`
		rows, err = m.DB.QueryContext(ctx, query, username)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var bookings []*models.SubmittedBooking
	for rows.Next() {
		var booking models.SubmittedBooking
		err := rows.Scan(
			&booking.ID,
			&booking.Username,
			&booking.Name,
			&booking.StartDate,
			&booking.EndDate,
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

func (m *PostgresDBRepo) InsertBookingRequest(booking models.Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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
		return err
	}

	// If overlap found, return error
	if err != sql.ErrNoRows {
		return errors.New("Booking time overlaps with an existing booking")
	}

	// If no overlaps, proceed with insertion
	stmt := `insert into requestedbookings (username, name, start_date, end_date, unit_number, start_time,
		end_time, purpose, facility)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var newID int

	err = m.DB.QueryRowContext(ctx, stmt,
		booking.Username,
		booking.Name,
		booking.StartDate,
		booking.EndDate,
		booking.UnitNumber,
		booking.StartTime,
		booking.EndTime,
		booking.Purpose,
		booking.Facility,
	).Scan(&newID)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) ApproveBookingRequest(booking models.SubmittedBooking) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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
		return err
	}

	// If overlap found, return error
	if err != sql.ErrNoRows {
		return errors.New("Booking time overlaps with an existing booking")
	}
	// If no overlaps, proceed with insertion
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Copy the booking from requestedbookings to approvedbookings
	copyStmt := `INSERT INTO approvedbookings (username, name, start_date, end_date, unit_number, start_time, end_time, purpose, facility)
		SELECT username, name, start_date, end_date, unit_number, start_time, end_time, purpose, facility FROM requestedbookings WHERE id = $1`
	_, err = tx.ExecContext(ctx, copyStmt, booking.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Delete the booking from requestedbookings
	deleteStmt := `DELETE FROM requestedbookings WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteStmt, booking.ID)
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

func (m *PostgresDBRepo) DeleteBookingRequest(booking models.SubmittedBooking) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Delete the booking from requestedbookings
	deleteStmt := `DELETE FROM requestedbookings WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteStmt, booking.ID)
	if err != nil {
		tx.Rollback() // Rollback in case of any error during the delete operation
		return err
	}

	err = tx.Commit() // Commit the transaction
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteApprovedBooking(booking models.SubmittedBooking) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Delete the booking from requestedbookings
	deleteStmt := `DELETE FROM approvedbookings WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteStmt, booking.ID)
	if err != nil {
		tx.Rollback() // Rollback in case of any error during the delete operation
		return err
	}

	err = tx.Commit() // Commit the transaction
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetUserByName(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, username, password, is_admin from users where username = $1`
	var user models.User
	row := m.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) RegisterUser(username string, password string, admin bool) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	stmt := `insert into users (username, password, is_admin) values ($1, $2, $3) returning id`

	var newID int

	err = m.DB.QueryRowContext(ctx, stmt, username, string(hashedPassword), admin).Scan(&newID)

	if err != nil {
		return nil, err
	}

	var user models.User = models.User{
		ID:       newID,
		Username: username,
		Password: password,
		IsAdmin:  admin,
	}

	return &user, nil
}
