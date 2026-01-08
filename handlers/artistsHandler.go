package handlers

import (
	"bytes"
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
)

// artistsHandler displays the main page with all artists
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	// only allow root path
	if r.URL.Path != "/" {
		renderError(w, "Page not found", http.StatusNotFound)
		return
	}

	// only allow GET requests
	if r.Method != http.MethodGet {
		renderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// fetch artists from API
	artists, err := api.GetArtists()
	if err != nil {
		log.Println(err)
		//if artist is not found
		if err.Error() == "artist not found." {
			renderError(w, "Artist not found.", http.StatusNotFound)
			return
		}
		renderError(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	// parse template
	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// execute template into buffer first
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, artists)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// send successful response
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}
