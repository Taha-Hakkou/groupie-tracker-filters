package utils

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/structures"
	"net/http"
	"slices"
	"strings"
)

// extractEvents populates artist with event data from multiple API endpoints
func ExtractEvents(artist structures.Artist) (structures.Artist, error) {
	// fetch location data
	resp, err := http.Get(artist.LocationsApi)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to fetch locations.")
	}

	if resp.StatusCode != http.StatusOK {
		return structures.Artist{}, fmt.Errorf("locations bad status code.")
	}

	decoder := json.NewDecoder(resp.Body)
	locationObject := structures.LocationObject{}
	err = decoder.Decode(&locationObject)
	resp.Body.Close()
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to decode locations.")
	}

	formatLocations(locationObject.Locations)

	// fetch date data
	resp, err = http.Get(artist.DatesApi)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to fetch dates.")
	}

	if resp.StatusCode != http.StatusOK {
		return structures.Artist{}, fmt.Errorf("dates bad status code.")
	}

	decoder = json.NewDecoder(resp.Body)
	dateObject := structures.DateObject{}
	err = decoder.Decode(&dateObject)
	resp.Body.Close()
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to decode dates.")
	}

	formatDates(dateObject.Dates)

	// fetch relation data (location->dates mapping)
	resp, err = http.Get(artist.RelationApi)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to fetch relations.")
	}

	if resp.StatusCode != http.StatusOK {
		return structures.Artist{}, fmt.Errorf("relations bad status code.")
	}

	decoder = json.NewDecoder(resp.Body)
	relationObject := structures.RelationObject{}
	err = decoder.Decode(&relationObject)
	resp.Body.Close()
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to decode relations.")
	}

	// build events from relations, validating against actual locations and dates
	for location, dates := range relationObject.LocationsDates {
		location = formatLocation(location)
		formatDates(dates)

		// skip if location not in valid list
		if !slices.Contains(locationObject.Locations, location) {
			continue
		}

		event := structures.Event{Location: location}

		// filter dates to only include valid ones
		for _, date := range dates {
			if slices.Contains(dateObject.Dates, date) {
				event.Dates = append(event.Dates, date)
			}
		}

		// only add event if it has dates
		if len(event.Dates) == 0 {
			continue
		}

		artist.Events = append(artist.Events, event)
	}

	return artist, nil
}

// formatDate removes leading asterisk from dates
func formatDate(date string) string {
	return strings.TrimPrefix(date, "*")
}

// formatDates formats a slice of dates in-place
func formatDates(dates []string) {
	for i := range dates {
		dates[i] = formatDate(dates[i])
	}
}

// formatLocation replaces dashes and underscores with spaces
func formatLocation(location string) string {
	location = strings.ReplaceAll(location, "-", " ")
	location = strings.ReplaceAll(location, "_", " ")
	return location
}

// formatLocations formats a slice of locations in-place
func formatLocations(locations []string) {
	for i := range locations {
		locations[i] = formatLocation(locations[i])
	}
}
