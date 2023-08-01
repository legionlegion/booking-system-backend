package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// variables of type application have access to this function
func (app *application) routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	// middleware and what routes we will have
	// applies to all requests to application

	// application logs when panic with backtraces
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Post("/authenticate", app.authenticate)
	mux.Post("/register", app.register)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authCheck)

		mux.Put("/add-booking", app.InsertBooking)
		mux.Put("/approve-booking", app.ApproveBooking)
		mux.Get("/booking-management", app.BookingManagement)
		mux.Put("/delete-pending", app.DeletePending)
		mux.Put("/delete-approved", app.DeleteApproved)

	})

	return mux
}
