package gtapi

import (
	"slices"
	"time"
)

type Range struct {
	From, To int
}

type TimeRange struct {
	From, To time.Time
}

type Filters struct {
	CreationYear   Range
	FirstAlbumDate TimeRange
	BandSizes      []int // number of members
	Country, City  string
}

func NewFilters(creationYear Range, firstAlbumYear TimeRange, bandsizes []int) Filters {
	filters := Filters{
		CreationYear:   creationYear,
		FirstAlbumDate: firstAlbumYear,
		BandSizes:      bandsizes,
	}
	// if not set return all artists !!!!!!!

	return filters
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
		filteredArtists = append(filteredArtists, artist)
	}
	return filteredArtists
}
