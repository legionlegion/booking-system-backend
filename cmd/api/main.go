package main

import (
	"booking-backend/internal/repository"
	"booking-backend/internal/repository/dbrepo"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	Domain       string
	DSN          string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

func main() {
	// set application config
	var app application

	// read from command line
	dsn, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		log.Print("DATABASE_URL not set, setting to localhost")
		flag.StringVar(
			&app.DSN,
			"dsn",
			"host=localhost port = 5432 user=syal password=syal dbname=bookings sslmode=disable timezone=UTC connect_timeout=5",
			"Postgres connection string",
		)
		flag.Parse()
	} else {
		app.DSN = dsn
	}

	flag.StringVar(&app.JWTSecret, "jwt-secret", "secret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "https://syal-2ae9b.firebaseapp.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "https://syal-2ae9b.firebaseapp.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "bookingsyal-cbd544b30b67.herokuapp.com", "cookie domain")
	// flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "https://syal-2ae9b.firebaseapp.com", "domain")
	flag.Parse()

	// connect to database
	conn, err := app.connectToDB()

	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close() // closes when main finishes running

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	port, exists := os.LookupEnv("PORT")

	if !exists {
		log.Print("Setting port to local")
		port = "8080"
	}

	// start a web server
	err = http.ListenAndServe(":"+port, app.routes()) // go-chi mux
	if err != nil {
		log.Fatal(err)
	}
}
