package gtapi

import (
	"slices"
	"strconv"
	"strings"
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
	filteredArtists := []Artist{}
	for _, artist := range artists {
		// creation year
		if artist.CreationDate < filters.CreationYear.From || artist.CreationDate > filters.CreationYear.To {
			continue
		}
		// first album year
		firstAlbumDate, _ := time.Parse("02-01-2006", artist.FirstAlbum)
		if firstAlbumDate.Before(filters.FirstAlbumDate.From) || firstAlbumDate.After(filters.FirstAlbumDate.To) {
			continue
		}
		// number of members
		if len(filters.BandSizes) > 0 && !slices.Contains(filters.BandSizes, len(artist.Members)) {
			continue
		}
		// location
		artist, _ = GetArtistDetails(strconv.Itoa(artist.Id))
		var countries []string
		var locations = map[string][]string{} // maps countries to cities
		for _, event := range artist.Events {
			fields := strings.Split(event.Location, "-")
			city, country := fields[0], fields[1]
			countries = append(countries, country)
			_, ok := locations[country]
			if ok {
				locations[country] = append(locations[country], city)
			} else {
				locations[country] = []string{city}
			}
		}
		if filters.Country != "" {
			if !slices.Contains(countries, strings.ToLower(filters.Country)) {
				continue
			}
			if filters.City != "" && !slices.Contains(locations[filters.Country], strings.ToLower(filters.City)) {
				continue
			}
		}

		filteredArtists = append(filteredArtists, artist)
	}
	return filteredArtists
}
