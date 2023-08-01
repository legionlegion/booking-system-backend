package repository

import (
	"booking-backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllBookings() ([]*models.SubmittedBooking, error)
	UserBookings(username string) ([]*models.SubmittedBooking, []*models.SubmittedBooking, []*models.RequestedBooking, error)
	AdminBookings() ([]*models.SubmittedBooking, []*models.SubmittedBooking, []*models.RequestedBooking, error)
	TwoWeekBookings() ([]*models.SubmittedBooking, error)
	ManageBookings(username string) ([]*models.SubmittedBooking, []*models.SubmittedBooking, []*models.RequestedBooking, error)
	InsertBookingRequest(booking models.Booking) error
	ApproveBookingRequest(booking models.RequestedBooking) error
	ApproveRecurringBookingRequest(booking models.RequestedBooking) error
	DeleteBookingRequest(booking models.SubmittedBooking) error
	DeleteApprovedBooking(booking models.SubmittedBooking) error
	GetUserByName(username string) (*models.User, error)
	RegisterUser(username string, password string, admin bool) (*models.User, error)
}
