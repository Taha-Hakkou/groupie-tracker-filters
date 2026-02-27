package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"groupie-tracker/gtapi"
)

// Displays the main page with all artists
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	// only allow root and artists path
	if !slices.Contains([]string{"/", "/artists", "/artists/"}, r.URL.Path) {
		renderError(w, "Page not found", http.StatusNotFound)
		return
	}

	// only allow GET requests
	if r.Method != http.MethodGet {
		renderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// fetch artists from API
	artists, err := gtapi.GetArtists()
	if err != nil {
		log.Println(err)
		renderError(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	// filter query params
	var errorMessage string

	// creation
	a := strings.TrimSpace(r.URL.Query().Get("cfrom"))
	b := strings.TrimSpace(r.URL.Query().Get("cto"))
	cfrom, err1 := strconv.Atoi(a) // default: 0
	cto, err2 := strconv.Atoi(b)
	if (a != "" && err1 != nil) || (b != "" && err2 != nil) {
		errorMessage = "creation year error: non-integer value"
	}
	if b == "" || err2 != nil {
		cto = time.Now().Year() // default: current year
	}
	var creationYear = gtapi.Range{
		From: cfrom,
		To:   cto,
	}

	// first album date
	layout := "2006-01-02"
	x := r.URL.Query().Get("afrom")
	y := r.URL.Query().Get("ato")
	from, err1 := time.Parse(layout, x)
	to, err2 := time.Parse(layout, y)
	if (x != "" && err1 != nil) || (y != "" && err2 != nil) {
		errorMessage = "first album year error: incorrect date format"
	}
	if y == "" || err2 != nil {
		to = time.Now().UTC() // default: current date
	}
	var firstAlbumYear = gtapi.TimeRange{
		From: from,
		To:   to,
	}

	// members
	members := r.URL.Query()["members"]
	var bandsizes []int
	for _, member := range members {
		n, err := strconv.Atoi(member)
		if err != nil {
			continue
		}
		bandsizes = append(bandsizes, n)
	}

	// location
	country := strings.TrimSpace(r.URL.Query().Get("country"))
	city := strings.TrimSpace(r.URL.Query().Get("city"))
	if country == "" && city != "" {
		errorMessage = "location error: country not specified"
	}

	var filters = gtapi.NewFilters(creationYear, firstAlbumYear, bandsizes, country, city)
	filteredArtists := gtapi.Filter(artists, filters)
	// -------------------------------------------

	// parse template
	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		log.Println("Error parsing artists template.")
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// execute template into buffer first
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, struct {
		Error   string
		Artists []gtapi.Artist
	}{
		Error:   errorMessage,
		Artists: filteredArtists,
	}) // artists
	if err != nil {
		log.Println("Error executing artists template.")
		renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// send successful response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}
