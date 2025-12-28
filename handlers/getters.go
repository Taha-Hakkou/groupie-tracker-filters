package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"groupie-tracker/models"
)

func getArtists() []models.Artist {
	res, err := http.Get(fmt.Sprintf("%s/artists", API))
	if err != nil {
		log.Fatal(err)
		// exits the server ?!
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	// var artists []interface{}
	var artists []models.Artist
	err = decoder.Decode(&artists)
	return artists
}

func getArtist(id int) models.Artist {
	res, err := http.Get(fmt.Sprintf("%s/artists/%d", API, id))
	if err != nil {
		log.Fatal(err)
		// exits the server ?!
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	// var artists []interface{}
	var artist models.Artist
	err = decoder.Decode(&artist)
	return artist
}

func getLocations() []models.Location {
	res, err := http.Get(fmt.Sprintf("%s/locations", API))
	if err != nil {
		log.Fatal(err)
		// exits the server ?!
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	// var artists []interface{}
	var locations []models.Location
	err = decoder.Decode(&locations)
	return locations
}

func getLocation(api string) models.Location {
	res, err := http.Get(fmt.Sprintf(api))
	if err != nil {
		log.Fatal(err)
		// exits the server ?!
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	// var artists []interface{}
	var location models.Location
	err = decoder.Decode(&location)
	return location
}

func getRelations() []models.Relation {
	res, err := http.Get(fmt.Sprintf("%s/relation", API))
	if err != nil {
		log.Fatal(err)
		// exits the server ?!
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	// var artists []interface{}
	var relations []models.Relation
	err = decoder.Decode(&relations)
	return relations
}
