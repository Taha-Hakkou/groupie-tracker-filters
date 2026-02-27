package gtapi

import "time"

type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	LocationsApi string  `json:"locations"`
	DatesApi     string  `json:"concertDates"`
	RelationApi  string  `json:"relations"`
	Events       []Event `json:"-"`
}

type Event struct {
	Location string
	Dates    []string
}

type LocationObject struct {
	Locations []string
}

type DateObject struct {
	Dates []string
}

type RelationObject struct {
	LocationsDates map[string][]string `json:"datesLocations"`
}

// ******** Filter structures **********

type Range struct {
	From, To int
}

type TimeRange struct {
	From, To time.Time
}

type Filters struct {
	CreationYear   Range
	FirstAlbumDate TimeRange
	BandSizes      []int
	Country, City  string
}
