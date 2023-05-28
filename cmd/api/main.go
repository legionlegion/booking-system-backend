package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 8080

type application struct {
	Domain string
}

func main() {
	// set application config
	var app application

	// read from command line

	// connect to database
	app.Domain = "example.com"

	log.Println("Running on port: ", port)

	// start a web server
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes()) // go-chi mux
	if err != nil {
		log.Fatal(err)
	}
}