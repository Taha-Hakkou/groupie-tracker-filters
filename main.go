package main

import (
	"log"
	"net/http"

	"groupie-tracker/handlers"
)

func main() {
	// register routes
	http.HandleFunc("/style.css", handlers.Styles)
	http.HandleFunc("/", handlers.Artists)
	http.HandleFunc("/artists/{id}", handlers.Artist)

	// start server
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
