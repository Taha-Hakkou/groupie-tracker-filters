package handlers

import (
	"fmt"
	"groupie-tracker/api"
	"html/template"
	"net/http"
)

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	//validate path
	if r.URL.Path != "/" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	//validate method
	if r.Method != http.MethodGet {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//get artists from api
	artists := api.GetArtists()
	//parse template
	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	//execute template
	err = tmpl.Execute(w, artists)
	if err != nil {
		fmt.Println(err)
	}

}
