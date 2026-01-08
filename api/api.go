package api

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/structures"
	"groupie-tracker/utils"
	"net/http"
)

const api = "https://groupietrackers.herokuapp.com/api/artists"

// getArtists fetches all artists from the API
func GetArtists() ([]structures.Artist, error) {
	resp, err := http.Get(api)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artists.")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artists bad status code.")
	}

	decoder := json.NewDecoder(resp.Body)
	artists := []structures.Artist{}
	err = decoder.Decode(&artists)
	if err != nil {
		return nil, fmt.Errorf("failed to decode artists.")
	}

	return artists, nil
}

// getArtistDetails fetches a single artist with full event details
func GetArtistDetails(id string) (structures.Artist, error) {
	artistEndpoint := api + "/" + id

	resp, err := http.Get(artistEndpoint)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to fetch artist.")
	}
	defer resp.Body.Close()

	// check for 404 or invalid artist ID
	if resp.StatusCode == http.StatusNotFound {
		return structures.Artist{}, fmt.Errorf("artist not found")
	}

	if resp.StatusCode != http.StatusOK {
		return structures.Artist{}, fmt.Errorf("failed to fetch artist.")
	}

	decoder := json.NewDecoder(resp.Body)
	artist := structures.Artist{}
	err = decoder.Decode(&artist)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to decode artist.")
	}

	// populate events with locations and dates
	populatedArtist, err := utils.ExtractEvents(artist)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to extract events... %w", err)
	}

	return populatedArtist, nil
}
