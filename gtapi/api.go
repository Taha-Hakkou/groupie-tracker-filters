package gtapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const api = "https://groupietrackers.herokuapp.com/api/artists"

// Fetches all artists from the API
func GetArtists() ([]Artist, error) {
	resp, err := http.Get(api)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artists")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artists bad status code")
	}

	decoder := json.NewDecoder(resp.Body)
	artists := []Artist{}
	err = decoder.Decode(&artists)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to decode artists")
	}

	return artists, nil
}

// Fetches a single artist with full event details
func GetArtistDetails(id string) (Artist, error) {
	// check if id is between 1 and 52
	idNumber, err := strconv.Atoi(id)
	if err != nil || (idNumber < 1 || idNumber > 52) {
		return Artist{}, fmt.Errorf("artist not found %s", id)
	}
	artistEndpoint := api + "/" + id

	resp, err := http.Get(artistEndpoint)
	if err != nil {
		return Artist{}, fmt.Errorf("failed to fetch artist %s", id)
	}

	// check for 404 or invalid artist ID
	if resp.StatusCode == http.StatusNotFound {
		return Artist{}, fmt.Errorf("artist not found %s", id)
	}

	if resp.StatusCode != http.StatusOK {
		return Artist{}, fmt.Errorf("failed to fetch artist %s", id)
	}

	decoder := json.NewDecoder(resp.Body)
	artist := Artist{}
	err = decoder.Decode(&artist)
	resp.Body.Close()
	if err != nil {
		return Artist{}, fmt.Errorf("failed to decode artist %s", id)
	}
	// populate events with locations and dates
	populatedArtist, err := ExtractEvents(artist)
	if err != nil {
		return Artist{}, fmt.Errorf("failed to extract artist %s events... %w", id, err)
	}

	return populatedArtist, nil
}
