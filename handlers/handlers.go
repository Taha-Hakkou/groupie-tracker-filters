package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"groupie-tracker/models"
)

const API string = "https://groupietrackers.herokuapp.com/api"

// Global CSS Handler
func CssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/global.css")
}

// All Artists Pahe Handler
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists := getArtists()
	t, _ := template.ParseFiles("templates/artists.html")
	t.Execute(w, struct {
		Artists []models.Artist
	}{
		Artists: artists,
	})
	// fmt.Fprintf(w, "%s", artists[1])
}

// Single Artist Page Handler
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	stringId := r.PathValue("id")
	id, _ := strconv.Atoi(stringId)
	artist := getArtist(id)

	artist.Locations = getLocation(artist.LocationsApi)

	t, _ := template.ParseFiles("templates/artist.html")
	t.Execute(w, artist)
	// fmt.Fprintf(w, "%s", artists[1])
}
