package main

import (
	"booking-backend/internal/repository"
	"booking-backend/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const port = 8080

type application struct {
	Domain string
	DSN    string
	DB     repository.DatabaseRepo
}

func main() {
	// set application config
	var app application

	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port = 5432 user=syal password=syal dbname=bookings sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()

	// connect to database
	app.Domain = "example.com"
	conn, err := app.connectToDB()

	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close() // closes when main finishes running

	log.Println("Running on port: ", port)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes()) // go-chi mux
	if err != nil {
		log.Fatal(err)
	}
}
