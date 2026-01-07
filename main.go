package main

import (
	"groupie-tracker/handlers"
	"log"
	"net/http"
)

func main() {
	// Register routes
	http.HandleFunc("/style.css", handlers.CssHandler)
	http.HandleFunc("/", handlers.ArtistsHandler)
	http.HandleFunc("/{id}", handlers.ArtistHandler)

	// Start server
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}