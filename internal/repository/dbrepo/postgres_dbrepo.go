package dbrepo

import (
	"booking-backend/internal/models"
	"context"
	"database/sql"
	"time"
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
