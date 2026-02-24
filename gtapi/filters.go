package gtapi

import "slices"

type Range struct {
	From, To int
}

type Filters struct {
	CreationYear, FirstAlbumYear Range
	BandSizes                    []int // number of members
	Country, City                string
}

func NewFilters(creationYear Range, bandsizes []int) Filters {
	filters := Filters{
		CreationYear: creationYear,
		BandSizes:    bandsizes,
	}
	// if not set return all artists !!!!!!!

	return filters
}

func Filter(artists []Artist, filters Filters) []Artist {
	filteredArtists := []Artist{}
	for _, artist := range artists {
		// creation
		if artist.CreationDate < filters.CreationYear.From || artist.CreationDate > filters.CreationYear.To {
			continue
		}
		// members
		if len(filters.BandSizes) > 0 && !slices.Contains(filters.BandSizes, len(artist.Members)) {
			continue
		}
		filteredArtists = append(filteredArtists, artist)
	}
	return filteredArtists
}
