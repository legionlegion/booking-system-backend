package models

import "time"

type Booking struct {
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	UnitNumber     string    `json:"unit_number"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	Purpose        string    `json:"purpose"`
	Facility       string    `json:"facility"`
	Recurring      bool      `json:"recurring"`
	RecurringWeeks int       `json:"recurring_weeks"`
}

type RequestedBooking struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	UnitNumber     string    `json:"unit_number"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	Purpose        string    `json:"purpose"`
	Facility       string    `json:"facility"`
	Recurring      bool      `json:"recurring"`
	RecurringWeeks int       `json:"recurring_weeks"`
}

type SubmittedBooking struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	UnitNumber string    `json:"unit_number"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	Purpose    string    `json:"purpose"`
	Facility   string    `json:"facility"`
}
