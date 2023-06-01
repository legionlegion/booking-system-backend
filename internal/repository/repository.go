package repository

import (
	"booking-backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllBookings() ([]*models.Booking, error)
	InsertBookingRequest(booking models.Booking) (int, error)
}
