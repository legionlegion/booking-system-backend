package main

import (
	"log"
	"net/http"
)

func (app *application) enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := map[string]bool{
			"https://syal-2ae9b.firebaseapp.com": true,
			"https://syal-2ae9b.web.app":         true,
			"http://localhost:3000":               true,
		}
		origin := r.Header.Get("Origin")
		log.Print("Origin: ", origin);
		log.Print("Allowed: ", allowedOrigins[origin]);
		if _, ok := allowedOrigins[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization, Username")
		}

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (app *application) authCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := app.auth.GetAndVerifyHeaderToken(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
