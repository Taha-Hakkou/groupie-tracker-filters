package api

import (
	"encoding/json"
	"groupie-tracker/structures"
	"log"
	"net/http"
)

func GetArtists() []structures.Artist {
	const api = "https://groupietrackers.herokuapp.com/api/artists"
	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(resp.Body) ///
	defer resp.Body.Close()               ///
	artists := []structures.Artist{}
	err = decoder.Decode(&artists) ///
	if err != nil {
		log.Fatal(err)
	}
	return artists
}
