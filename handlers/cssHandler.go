package handlers

import (
	"net/http"
)

// CssHandler serves the CSS file
func CssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/style.css")
}