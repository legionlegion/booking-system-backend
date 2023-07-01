package repository

import (
	"booking-backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllBookings() ([]*models.SubmittedBooking, error)
	ManageBookings(username string) ([]*models.SubmittedBooking, error)
	InsertBookingRequest(booking models.Booking) error
	ApproveBookingRequest(booking models.SubmittedBooking) error
	DeleteBookingRequest(booking models.SubmittedBooking) error
	GetUserByName(username string) (*models.User, error)
	RegisterUser(username, password string) (*models.User, error)
}
