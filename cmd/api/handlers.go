package main

import (
	"booking-backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// all handlers take 2 arguements
// 1: ResponseWriter writes response to client
// 2: pointer to a http request (not actual value but memory address)
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello world! This is %s", app.Domain)
	// var payload = struct {
	// 	Status string `json:"status"`
	// 	Message string `json:"message"`
	// 	Version string `json:"version"`
	// }{
	// 	Status: "active",
	// 	Message: "Go Movies up and running",
	// 	Version: "1.0.0",
	// }

	var bookings []models.Booking

	booking1 := models.Booking {
		Date: time.Now(),
		StartTime: 8,
		EndTime: 10,
	}

	bookings = append(bookings, booking1)

	// out, err := json.Marshal(payload)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	out, err := json.Marshal(booking1)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
} 