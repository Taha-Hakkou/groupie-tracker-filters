package handlers

import (
	"bytes"
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
	"strings"
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
		log.Println(err)
		//if artist is not found
		if strings.HasPrefix(err.Error(), "artist not found") {
			renderError(w, "Artist not found.", http.StatusNotFound)
			return
		}
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// parse template
	tmpl, err := template.ParseFiles("templates/artist-details.html")
	if err != nil {
		log.Println("Error parsing artist-details template.")
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// execute template into buffer first
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, artist)
	if err != nil {
		log.Println("Error executing artist-details template.")
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// send successful response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}
