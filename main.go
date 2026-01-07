package main

import (
	"fmt"
	"groupie-tracker/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/style.css", handlers.CssHandler)
	http.HandleFunc("/", handlers.ArtistsHandler)
	http.HandleFunc("/{id}", handlers.ArtistHandler)
	fmt.Println("Server listening on localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("could not start server... %s", err)
	}
}
