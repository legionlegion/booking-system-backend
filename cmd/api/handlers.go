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
	log.Print("BACKEND HOME")
	bookings, err := app.DB.AllBookings()
	if err != nil {
		log.Print("HOME ERR: ", err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, bookings)
}

func (app *application) InsertBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	log.Print("Booking in handler: ", booking)

	err := app.readJSON(w, r, &booking)
	if err != nil {
		log.Print("Error reading json: ", err)
		app.errorJSON(w, err)
		return
	}
	err = app.DB.InsertBookingRequest(booking)
	log.Print("Err in insertbooking: ", err)
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

	log.Println(tokens.Token)
	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
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

	log.Print("Register err 1: ", err)
	// register user
	user, err := app.DB.RegisterUser(requestPayload.Username, requestPayload.Password)

	// create a jwt user
	u := jwtUser{
		ID:       user.ID,
		Username: user.Username,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Println(tokens.Token)
	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the username from the token claims
			username := claims.Username
			if err != nil {

				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			// check if user still exists
			user, err := app.DB.GetUserByName(username)

			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:       user.ID,
				Username: user.Username,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)

			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

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
