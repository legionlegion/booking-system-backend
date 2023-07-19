package main

import (
	"booking-backend/internal/models"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
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
	bookings, err := app.DB.TwoWeekBookings()
	if err != nil {
		log.Println("Err: ", err)
		return
	}
	log.Print("Bookings: ")

	_ = app.writeJSON(w, http.StatusOK, bookings)
}

func (app *application) AllBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := app.DB.AllBookings()
	if err != nil {
		return
	}

	_ = app.writeJSON(w, http.StatusOK, bookings)
}

func (app *application) InsertBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking

	err := app.readJSON(w, r, &booking)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.DB.InsertBookingRequest(booking)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Booking requested",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) ApproveBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.SubmittedBooking

	err := app.readJSON(w, r, &booking)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.DB.ApproveBookingRequest(booking)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Booking requested",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) DeletePending(w http.ResponseWriter, r *http.Request) {
	var booking models.SubmittedBooking

	err := app.readJSON(w, r, &booking)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.DB.DeleteBookingRequest(booking)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Booking requested",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) DeleteApproved(w http.ResponseWriter, r *http.Request) {
	var booking models.SubmittedBooking

	err := app.readJSON(w, r, &booking)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.DB.DeleteApprovedBooking(booking)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Booking requested",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	// validate user against database
	user, err := app.DB.GetUserByName(requestPayload.Username)
	if err != nil {
		app.errorJSON(w, errors.New("Username does not exist"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Admin    bool   `json:"admin"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// register user
	user, err := app.DB.RegisterUser(requestPayload.Username, requestPayload.Password, requestPayload.Admin)

	// create a jwt user
	u := jwtUser{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		// log.Println("Request Host: ", r.Host)
		// log.Println("Request Path: ", r.URL.Path)
		// log.Println("Cookie name: ", cookie.Name)
		// log.Println("App cookie name: ", app.auth.CookieName)
		if cookie.Name == app.auth.CookieName {
			log.Println("Cookie name passed")
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				log.Println("Success claim parsing")
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				log.Println("Unauthorized")
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the username from the token claims
			username := claims.Username
			if err != nil {
				log.Println("Unknown user")
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			// check if user still exists
			user, err := app.DB.GetUserByName(username)

			if err != nil {
				log.Println("Unknown user")
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:       user.ID,
				Username: user.Username,
				IsAdmin:  user.IsAdmin,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)

			if err != nil {
				log.Println("Error generating token")
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			log.Println("Sucess token generation")
			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPairs)
		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) BookingManagement(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("Username")
	bookings, err := app.DB.ManageBookings(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, bookings)
}

func (app *application) UserBookings(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("Username")
	bookings, err := app.DB.UserBookings(username)
	if err != nil {
		return
	}

	_ = app.writeJSON(w, http.StatusOK, bookings)
}