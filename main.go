package main

import (
	"log"
	"net/http"

	"groupie-tracker/handlers"
)

func main() {
	// register routes
	http.HandleFunc("/style.css", handlers.CssHandler)
	http.HandleFunc("/", handlers.ArtistsHandler)
	http.HandleFunc("/artists/{id}", handlers.ArtistHandler)

	// start server
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
