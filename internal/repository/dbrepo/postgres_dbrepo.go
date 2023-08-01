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

func (m *PostgresDBRepo) AllBookings() ([]*models.SubmittedBooking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select 
		id, username, name, start_date, end_date, unit_number, 
		start_time, end_time, purpose, facility 
	from 
		approvedbookings 
	order by 
		start_date ASC, start_time ASC
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
		start_date >= date_trunc('week', current_date) 
		and end_date < date_trunc('week', current_date) + interval '2 weeks'
	order by 
		start_date ASC, start_time ASC
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
	log.Println("Bookings: ", bookings)
	return bookings, nil
}

func (m *PostgresDBRepo) ManageBookings(username string) ([]*models.SubmittedBooking, []*models.SubmittedBooking, []*models.RequestedBooking, error) {
	user, err := m.GetUserByName(username)
	if err != nil {
		return nil, nil, nil, err
	}

	if user.IsAdmin {
		// get all bookings
		return m.AdminBookings()
	} else {
		// get bookings that belong to user only
		return m.UserBookings(username)
	}
}

func (m *PostgresDBRepo) AdminBookings() ([]*models.SubmittedBooking, []*models.SubmittedBooking, []*models.RequestedBooking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Query for recurring
	recurringQuery := `
		SELECT 
			id, username, name, start_date, end_date, unit_number, 
			start_time, end_time, purpose, facility 
		FROM 
			recurringbookings
		ORDER BY 
			start_date ASC, start_time ASC
	`
	recurringRows, err := m.DB.QueryContext(ctx, recurringQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	defer recurringRows.Close()

	// Query for approvedbookings
	approvedQuery := `
		SELECT 
			id, username, name, start_date, end_date, unit_number, 
			start_time, end_time, purpose, facility 
		FROM 
			approvedbookings
		ORDER BY 
			start_date ASC, start_time ASC
	`
	approvedRows, err := m.DB.QueryContext(ctx, approvedQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	defer approvedRows.Close()

	// Query for requestedbookings
	requestedQuery := `
		SELECT 
			id, username, name, start_date, end_date, unit_number, 
			start_time, end_time, purpose, facility, is_recurring, recurring_weeks 
		FROM 
			requestedbookings
		ORDER BY 
			start_date ASC, start_time ASC
	`
	requestedRows, err := m.DB.QueryContext(ctx, requestedQuery)
	if err != nil {
		return nil, nil, nil, err
	}
	defer requestedRows.Close()

	// Scan recurring bookings
	var recurringbookings []*models.SubmittedBooking
	for recurringRows.Next() {
		var booking models.SubmittedBooking
		err := recurringRows.Scan(
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
			return nil, nil, nil, err
		}
		recurringbookings = append(recurringbookings, &booking)
	}

	// Scan approved bookings
	var approvedBookings []*models.SubmittedBooking
	for approvedRows.Next() {
		var booking models.SubmittedBooking
		err := approvedRows.Scan(
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
			return nil, nil, nil, err
		}
		approvedBookings = append(approvedBookings, &booking)
	}

	// Scan requested bookings
	var requestedBookings []*models.RequestedBooking
	for requestedRows.Next() {
		var booking models.RequestedBooking
		err := requestedRows.Scan(
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
			&booking.Recurring,
			&booking.RecurringWeeks,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		requestedBookings = append(requestedBookings, &booking)
	}

	return recurringbookings, approvedBookings, requestedBookings, nil
}

func (m *PostgresDBRepo) UserBookings(username string) ([]*models.SubmittedBooking, []*models.SubmittedBooking, []*models.RequestedBooking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Query for recurring
	recurringQuery := `
		SELECT 
			id, username, name, start_date, end_date, unit_number, 
			start_time, end_time, purpose, facility 
		FROM 
			recurringbookings
		WHERE
			username = $1
		ORDER BY 
			start_date ASC, start_time ASC
	`
	recurringRows, err := m.DB.QueryContext(ctx, recurringQuery, username)
	if err != nil {
		return nil, nil, nil, err
	}
	defer recurringRows.Close()

	// Query for approvedbookings
	approvedQuery := `
		SELECT 
			id, username, name, start_date, end_date, unit_number, 
			start_time, end_time, purpose, facility 
		FROM 
			approvedbookings
		WHERE
			username = $1
		ORDER BY 
			start_date ASC, start_time ASC
	`
	approvedRows, err := m.DB.QueryContext(ctx, approvedQuery, username)
	if err != nil {
		return nil, nil, nil, err
	}
	defer approvedRows.Close()

	// Query for requestedbookings
	requestedQuery := `
		SELECT 
			id, username, name, start_date, end_date, unit_number, 
			start_time, end_time, purpose, facility, is_recurring, recurring_weeks 
		FROM 
			requestedbookings
		WHERE
			username = $1
		ORDER BY 
			start_date ASC, start_time ASC
	`
	requestedRows, err := m.DB.QueryContext(ctx, requestedQuery, username)
	if err != nil {
		return nil, nil, nil, err
	}
	defer requestedRows.Close()

	// Scan recurring bookings
	var recurringbookings []*models.SubmittedBooking
	for recurringRows.Next() {
		var booking models.SubmittedBooking
		err := recurringRows.Scan(
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
			return nil, nil, nil, err
		}
		recurringbookings = append(recurringbookings, &booking)
	}

	// Scan approved bookings
	var approvedBookings []*models.SubmittedBooking
	for approvedRows.Next() {
		var booking models.SubmittedBooking
		err := approvedRows.Scan(
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
			return nil, nil, nil, err
		}
		approvedBookings = append(approvedBookings, &booking)
	}

	// Scan requested bookings
	var requestedBookings []*models.RequestedBooking
	for requestedRows.Next() {
		var booking models.RequestedBooking
		err := requestedRows.Scan(
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
			&booking.Recurring,
			&booking.RecurringWeeks,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		requestedBookings = append(requestedBookings, &booking)
	}

	return recurringbookings, approvedBookings, requestedBookings, nil
}

func (m *PostgresDBRepo) InsertBookingRequest(booking models.Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// check for overlaps, WHERE clause covers all overlap scenarios
	checkOverlapStmt := `
		SELECT id
		FROM (
			SELECT id, start_time, end_time
			FROM approvedbookings
			UNION ALL
			SELECT id, start_time, end_time
			FROM recurringbookings
		) AS all_bookings
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
		return errors.New("booking time overlaps with an existing booking")
	}

	// If no overlaps, proceed with insertion
	stmt := `insert into requestedbookings (username, name, start_date, end_date, unit_number, start_time,
		end_time, purpose, facility, is_recurring, recurring_weeks)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id`

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
		booking.Recurring,
		booking.RecurringWeeks,
	).Scan(&newID)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) ApproveBookingRequest(booking models.RequestedBooking) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// check for overlaps, WHERE clause covers all overlap scenarios
	checkOverlapStmt := `
		SELECT id
		FROM (
			SELECT id, start_time, end_time
			FROM approvedbookings
			UNION ALL
			SELECT id, start_time, end_time
			FROM recurringbookings
		) AS all_bookings
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
		return errors.New("booking time overlaps with an existing booking")
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

func (m *PostgresDBRepo) ApproveRecurringBookingRequest(booking models.RequestedBooking) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// check for overlaps, WHERE clause covers all overlap scenarios
	checkOverlapStmt := `
		SELECT id
		FROM (
			SELECT id, start_time, end_time
			FROM approvedbookings
			UNION ALL
			SELECT id, start_time, end_time
			FROM recurringbookings
		) AS all_bookings
		WHERE 
			($1 < end_time AND $2 > start_time)
		LIMIT 1;
		`

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	layout := "2006-01-02T15:04:05-07:00"

	for week := 0; week < booking.RecurringWeeks; week++ {
		// Parse the start and end times to time.Time
		startTime, err := time.Parse(layout, booking.StartTime)
		if err != nil {
			log.Print("Error parsing start time: ", err)
			_ = tx.Rollback()
			return err
		}
		endTime, err := time.Parse(layout, booking.EndTime)
		if err != nil {
			log.Print("Error parsing end time: ", err)
			_ = tx.Rollback()
			return err
		}

		// Calculate the start and end time for this booking
		startTime = startTime.Add(time.Duration(week) * 7 * 24 * time.Hour)
		endTime = endTime.Add(time.Duration(week) * 7 * 24 * time.Hour)
		// Check for overlaps
		var overlapID int
		err = m.DB.QueryRowContext(ctx, checkOverlapStmt, startTime, endTime).Scan(&overlapID)
		if err != nil && err != sql.ErrNoRows {
			log.Print("Overlap check error: ", err)
			_ = tx.Rollback()
			return err
		}

		if err == sql.ErrNoRows {
			// If there is no overlap, insert the booking into the recurringbookings table
			insertStmt := `INSERT INTO recurringbookings (username, name, start_date, end_date, unit_number, start_time, end_time, purpose, facility)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
			_, err = tx.ExecContext(ctx, insertStmt, booking.Username, booking.Name, startTime, endTime, booking.UnitNumber, startTime, endTime, booking.Purpose, booking.Facility)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
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
