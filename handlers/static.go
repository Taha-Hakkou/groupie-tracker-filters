package handlers

import (
	"net/http"
)

// Serves the CSS file
func Styles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "assets/style.css")
}
