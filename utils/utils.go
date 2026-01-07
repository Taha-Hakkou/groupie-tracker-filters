package utils

import (
	"encoding/json"
	"groupie-tracker/structures"
	"net/http"
	"slices"
	"strings"
)

func ExtractEvents(artist structures.Artist) (structures.Artist, error) {
	//get locationObject
	resp, err := http.Get(artist.LocationsApi)
	if err != nil {
		return structures.Artist{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	locationObject := structures.LocationObject{}
	err = decoder.Decode(&locationObject)
	if err != nil {
		return structures.Artist{}, err
	}
	resp.Body.Close()
	formatLocations(locationObject.Locations)
	//get dateObject
	resp, err = http.Get(artist.DatesApi)
	if err != nil {
		return structures.Artist{}, err
	}
	decoder = json.NewDecoder(resp.Body)
	dateObject := structures.DateObject{}
	err = decoder.Decode(&dateObject)
	if err != nil {
		return structures.Artist{}, err
	}
	resp.Body.Close()
	formatDates(dateObject.Dates)
	//get relationObject
	resp, err = http.Get(artist.RelationApi)
	if err != nil {
		return structures.Artist{}, err
	}
	decoder = json.NewDecoder(resp.Body)
	relationObject := structures.RelationObject{}
	err = decoder.Decode(&relationObject)
	if err != nil {
		return structures.Artist{}, err
	}
	resp.Body.Close()
	//----------------------------------------
	//populate events
	for location, dates := range relationObject.LocationsDates {
		location = formatLocation(location)
		formatDates(dates)
		//match the relation location against the locations slice
		if !slices.Contains(locationObject.Locations, location) {
			continue
		}
		event := structures.Event{Location: location}
		for _, date := range dates {
			//match the relation location->date against the dates slice
			if !slices.Contains(dateObject.Dates, date) {
				continue
			}
			event.Dates = append(event.Dates, date)
		}
		//check if number of dates is not 0
		if len(event.Dates) == 0 {
			continue
		}
		artist.Events = append(artist.Events, event)
	}
	return artist, nil
}

func formatDate(date string) string {
	if strings.HasPrefix(date, "*") {
		dateRunes := []rune(date)
		newDateRunes := dateRunes[1:]
		return string(newDateRunes)
	}
	return date
}
func formatDates(dates []string) {
	for i := range dates {
		dates[i] = formatDate(dates[i])
	}
}
func formatLocation(location string) string {
	locationRunes := []rune(location)
	for i := range locationRunes {
		if locationRunes[i] == '-' || locationRunes[i] == '_' {
			locationRunes[i] = ' '
		}
	}
	return string(locationRunes)
}
func formatLocations(locations []string) {
	for i := range locations {
		locations[i] = formatLocation(locations[i])
	}
}
