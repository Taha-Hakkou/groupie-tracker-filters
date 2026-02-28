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
	var filteredArtists []Artist
	for _, artist := range artists {
		//creation date
		if artist.CreationDate < filters.CreationYear.From || artist.CreationDate > filters.CreationYear.To {
			continue
		}
		//first album
		firstAlbumDate, _ := time.Parse("02-01-2006", artist.FirstAlbum)
		if firstAlbumDate.Before(filters.FirstAlbumDate.From) || firstAlbumDate.After(filters.FirstAlbumDate.To) {
			continue
		}
		//number of members
		if len(filters.BandSizes) > 0 && !slices.Contains(filters.BandSizes, len(artist.Members)) {
			continue
		}
		filteredArtists = append(filteredArtists, artist)
	}

	// No location filter — return filtered artists as-is
	if filters.Country == "" {
		return filteredArtists
	}

	// Fetch details for all filtered artists in parallel
	type result struct {
		index  int
		artist Artist
		ok     bool
	}
	results := make([]result, len(filteredArtists))
	var wg sync.WaitGroup

	for i, artist := range filteredArtists {
		wg.Add(1)
		go func(i int, artist Artist) {
			defer wg.Done()
			detailed, err := GetArtistDetails(strconv.Itoa(artist.Id))
			if err != nil {
				results[i] = result{index: i, ok: false}
				return
			}

			formattedCountry := strings.ToLower(filters.Country)
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
				results[i] = result{index: i, ok: false}
				return
			}

			if filters.City != "" {
				formattedCity := strings.ToLower(filters.City)
				formattedCity = strings.ReplaceAll(formattedCity, " ", "_")
				if !slices.Contains(locations[formattedCountry], formattedCity) {
					results[i] = result{index: i, ok: false}
					return
				}
			}

			results[i] = result{index: i, artist: detailed, ok: true}
		}(i, artist)
	}

	wg.Wait()

	// Preserve order
	var filtered []Artist
	for _, r := range results {
		if r.ok {
			filtered = append(filtered, r.artist)
		}
	}
	return filtered
}