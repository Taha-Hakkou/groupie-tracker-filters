package handlers

import (
	"bytes"
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
)

// ArtistsHandler displays the main page with all artists
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow root path
	if r.URL.Path != "/" {
		renderError(w, "Page not found", http.StatusNotFound)
		return
	}

	// Only allow GET requests
	if r.Method != http.MethodGet {
		renderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch artists from API
	artists, err := api.GetArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		renderError(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	// Parse template
	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Execute template into buffer first
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, artists)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send successful response
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}