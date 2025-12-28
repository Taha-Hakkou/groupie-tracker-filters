package main

import (
	"log"
	"net/http"

	"groupie-tracker/handlers"
)

func main() {
	http.HandleFunc("/global.css", handlers.CssHandler)
	http.HandleFunc("/artists", handlers.ArtistsHandler)
	http.HandleFunc("/artists/{id}", handlers.ArtistHandler)

	log.Println("Server listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
