package utils

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/structures"
	"net/http"
	"slices"
	"strings"
)

// ExtractEvents populates artist with event data from multiple API endpoints
func ExtractEvents(artist structures.Artist) (structures.Artist, error) {
	// Fetch location data
	resp, err := http.Get(artist.LocationsApi)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return structures.Artist{}, fmt.Errorf("locations API returned status %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	locationObject := structures.LocationObject{}
	err = decoder.Decode(&locationObject)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to decode locations: %w", err)
	}

	formatLocations(locationObject.Locations)

	// Fetch date data
	resp, err = http.Get(artist.DatesApi)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to fetch dates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return structures.Artist{}, fmt.Errorf("dates API returned status %d", resp.StatusCode)
	}

	decoder = json.NewDecoder(resp.Body)
	dateObject := structures.DateObject{}
	err = decoder.Decode(&dateObject)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to decode dates: %w", err)
	}

	formatDates(dateObject.Dates)

	// Fetch relation data (location->dates mapping)
	resp, err = http.Get(artist.RelationApi)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to fetch relations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return structures.Artist{}, fmt.Errorf("relations API returned status %d", resp.StatusCode)
	}

	decoder = json.NewDecoder(resp.Body)
	relationObject := structures.RelationObject{}
	err = decoder.Decode(&relationObject)
	if err != nil {
		return structures.Artist{}, fmt.Errorf("failed to decode relations: %w", err)
	}

	// Build events from relations, validating against actual locations and dates
	for location, dates := range relationObject.LocationsDates {
		location = formatLocation(location)
		formatDates(dates)

		// Skip if location not in valid list
		if !slices.Contains(locationObject.Locations, location) {
			continue
		}

		event := structures.Event{Location: location}

		// Filter dates to only include valid ones
		for _, date := range dates {
			if slices.Contains(dateObject.Dates, date) {
				event.Dates = append(event.Dates, date)
			}
		}

		// Only add event if it has dates
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

// formatLocation replaces hyphens and underscores with spaces
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