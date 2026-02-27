package handlers

import (
	"bytes"
	"html/template"
	"net/http"
)

// Displays error page with proper status code
func renderError(w http.ResponseWriter, message string, statusCode int) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		// fallback if error template fails
		http.Error(w, message, statusCode)
		return
	}

	data := struct {
		Message    string
		StatusCode int
	}{
		Message:    message,
		StatusCode: statusCode,
	}

	// execute into buffer to catch any template errors
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		http.Error(w, message, statusCode)
		return
	}

	// send error response with proper status code
	w.WriteHeader(statusCode)
	buf.WriteTo(w)
}
