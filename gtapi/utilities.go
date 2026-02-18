package gtapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// Populates artist with event data from multiple API endpoints
func ExtractEvents(artist Artist) (Artist, error) {
	// fetch location data
	resp, err := http.Get(artist.LocationsApi)
	if err != nil {
		return Artist{}, fmt.Errorf("failed to fetch locations")
	}

	if resp.StatusCode != http.StatusOK {
		return Artist{}, fmt.Errorf("locations bad status code")
	}

	decoder := json.NewDecoder(resp.Body)
	locationObject := LocationObject{}
	err = decoder.Decode(&locationObject)
	resp.Body.Close()
	if err != nil {
		return Artist{}, fmt.Errorf("failed to decode locations")
	}

	formatLocations(locationObject.Locations)

	// fetch date data
	resp, err = http.Get(artist.DatesApi)
	if err != nil {
		return Artist{}, fmt.Errorf("failed to fetch dates")
	}

	if resp.StatusCode != http.StatusOK {
		return Artist{}, fmt.Errorf("dates bad status code")
	}

	decoder = json.NewDecoder(resp.Body)
	dateObject := DateObject{}
	err = decoder.Decode(&dateObject)
	resp.Body.Close()
	if err != nil {
		return Artist{}, fmt.Errorf("failed to decode dates")
	}

	formatDates(dateObject.Dates)

	// fetch relation data (location->dates mapping)
	resp, err = http.Get(artist.RelationApi)
	if err != nil {
		return Artist{}, fmt.Errorf("failed to fetch relations")
	}

	if resp.StatusCode != http.StatusOK {
		return Artist{}, fmt.Errorf("relations bad status code")
	}

	decoder = json.NewDecoder(resp.Body)
	relationObject := RelationObject{}
	err = decoder.Decode(&relationObject)
	resp.Body.Close()
	if err != nil {
		return Artist{}, fmt.Errorf("failed to decode relations")
	}

	// build events from relations, validating against actual locations and dates
	for location, dates := range relationObject.LocationsDates {
		location = formatLocation(location)
		formatDates(dates)

		// skip if location not in valid list
		if !slices.Contains(locationObject.Locations, location) {
			continue
		}

		event := Event{Location: location}

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

// Removes leading asterisk from dates
// func formatDate(date string) string {
// 	return strings.TrimPrefix(date, "*")
// }

// Formats a slice of dates in-place
func formatDates(dates []string) {
	for i := range dates {
		dates[i] = strings.TrimPrefix(dates[i], "*")
	}
}

// Replaces dashes and underscores with spaces
func formatLocation(location string) string {
	location = strings.ReplaceAll(location, "-", " ")
	location = strings.ReplaceAll(location, "_", " ")
	return location
}

// Formats a slice of locations in-place
func formatLocations(locations []string) {
	for i := range locations {
		locations[i] = formatLocation(locations[i])
	}
}
