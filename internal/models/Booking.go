package models

import "time"

type Booking struct {
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	UnitNumber string    `json:"unit_number"`
	Date       time.Time `json:"date"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	Purpose    string    `json:"purpose"`
	Facility   string    `json:"facility"`
}

type ApprovedBooking struct {
	ID      int       `json:"id"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	UnitNumber string    `json:"unit_number"`
	Date       time.Time `json:"date"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	Purpose    string    `json:"purpose"`
	Facility   string    `json:"facility"`
}
