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
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)
	return mux
}