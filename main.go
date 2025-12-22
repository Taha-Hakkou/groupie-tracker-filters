package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

const API string = "https://groupietrackers.herokuapp.com/api"

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(fmt.Sprintf("%s/artists", API))
	if err != nil {
		log.Fatal(err)
		// exits the server ?!
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	//var artists []interface{}
	var artists []Artist
	err = decoder.Decode(&artists)
	t, _ := template.ParseFiles("templates/artists.html")
	t.Execute(w, struct {
		Artists []Artist
	}{
		Artists: artists,
	})
	//fmt.Fprintf(w, "%s", artists[1])
}

func main() {
	http.HandleFunc("/artists", artistsHandler)
	log.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
