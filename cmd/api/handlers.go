package main

import (
	"booking-backend/internal/models"
	"fmt"
	"net/http"
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
	
	bookings, err := app.DB.AllBookings()
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, bookings)
} 

func (app *application) InsertBooking (w http.ResponseWriter, r *http.Request) {
	var booking models.Booking

	err := app.readJSON(w, r, &booking)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error: false,
		Message: "Booking requested",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}