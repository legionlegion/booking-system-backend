package models

import "time"

type Booking struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Name string `json:"name"`
	UnitNumber string `json:"unit_number"`
	Date time.Time `json:"date"`
	StartTime int `json:"start_time"`
	EndTime int `json:"end_time"`
	Purpose string `json:"purpose"`
	Facility string `json:"facility"`
}