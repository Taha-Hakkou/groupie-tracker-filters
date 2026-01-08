package handlers

import (
	"bytes"
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
)

// artistHandler displays individual artist details
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	// only allow GET requests
	if r.Method != http.MethodGet {
		renderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// extract artist ID from URL
	stringId := r.PathValue("id")

	// fetch artist details from API
	artist, err := api.GetArtistDetails(stringId)
	if err != nil {
		log.Printf("Error fetching artist %s: %v", stringId, err)
		renderError(w, "Artist not found", http.StatusNotFound)
		return
	}

	// parse template
	tmpl, err := template.ParseFiles("templates/artist-details.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// execute template into buffer first
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, artist)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// send successful response
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}
