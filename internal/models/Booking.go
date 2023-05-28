package models

import "time"

type Booking struct {
	Date time.Time `json:"date"`
	StartTime int `json:"start_time"`
	EndTime int `json:"end_time"`
}