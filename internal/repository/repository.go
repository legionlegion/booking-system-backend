package repository

import (
	"booking-backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllBookings() ([]*models.ApprovedBooking, error)
	InsertBookingRequest(booking models.Booking) (int, error)
	GetUserByName(username string) (*models.User, error)
	RegisterUser(username, password string) (*models.User, error)
}
