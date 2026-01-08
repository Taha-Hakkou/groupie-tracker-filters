package handlers

import (
	"net/http"
)

// cssHandler serves the CSS file
func CssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/style.css")
}
