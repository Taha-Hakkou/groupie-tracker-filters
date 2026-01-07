package handlers

import (
	"groupie-tracker/api"
	"net/http"
)

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	//validate method
	if r.Method != http.MethodGet {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//extract parameter
	stringId := r.PathValue("id")
	//get artist from api
	_, err := api.GetArtistDetails(stringId)
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

}
