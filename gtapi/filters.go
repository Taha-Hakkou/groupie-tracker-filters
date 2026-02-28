package gtapi

import (
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

func New(creationYear Range, firstAlbumYear TimeRange, bandsizes []int, country, city string) Filters {
	return Filters{
		CreationYear:   creationYear,
		FirstAlbumDate: firstAlbumYear,
		BandSizes:      bandsizes,
		Country:        country,
		City:           city,
	}
}

func Filter(artists []Artist, filters Filters) []Artist {
	// Apply cheap filters first (no HTTP calls)
	var semiFilteredArtists []Artist
	for _, artist := range artists {
		// creation date
		if artist.CreationDate < filters.CreationYear.From || artist.CreationDate > filters.CreationYear.To {
			continue
		}
		// first album
		firstAlbumDate, _ := time.Parse("02-01-2006", artist.FirstAlbum)
		if firstAlbumDate.Before(filters.FirstAlbumDate.From) || firstAlbumDate.After(filters.FirstAlbumDate.To) {
			continue
		}
		// number of members
		if len(filters.BandSizes) > 0 && !slices.Contains(filters.BandSizes, len(artist.Members)) {
			continue
		}
		semiFilteredArtists = append(semiFilteredArtists, artist)
	}

	// location
	// No location filter — return filtered artists as-is
	if filters.Country == "" {
		return semiFilteredArtists
	}

	var filteredArtists []Artist
	var wg sync.WaitGroup
	for _, artist := range semiFilteredArtists {
		wg.Add(1)
		go FilterByLocation(&wg, &filteredArtists, artist, filters.Country, filters.City)
	}
	wg.Wait()

	SortArtistsById(filteredArtists) // go-routines breaks the order
	return filteredArtists
}

func FilterByLocation(wg *sync.WaitGroup, artists *[]Artist, artist Artist, fcountry, fcity string) {
	defer wg.Done()
	detailed, err := GetArtistDetails(strconv.Itoa(artist.Id))
	if err != nil {
		return
	}

	formattedCountry := strings.ToLower(fcountry)
	var countries []string
	locations := map[string][]string{}
	for _, event := range detailed.Events {
		fields := strings.Split(event.Location, "-")
		if len(fields) < 2 {
			continue
		}
		city, country := fields[0], fields[1]
		countries = append(countries, country)
		locations[country] = append(locations[country], city)
	}

	if !slices.Contains(countries, formattedCountry) {
		return
	}

	if fcity != "" {
		formattedCity := strings.ToLower(fcity)
		formattedCity = strings.ReplaceAll(formattedCity, " ", "_")
		if !slices.Contains(locations[formattedCountry], formattedCity) {
			return
		}
	}
	*artists = append(*artists, artist)
}

func SortArtistsById(artists []Artist) {
	for i := range len(artists) {
		for j := range len(artists) - i - 1 {
			if artists[j].Id > artists[j+1].Id {
				artists[j], artists[j+1] = artists[j+1], artists[j]
			}
		}
	}
}
