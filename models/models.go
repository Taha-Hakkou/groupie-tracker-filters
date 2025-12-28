package models

type Artist struct {
	Id              int
	Image           string
	Name            string
	Members         []string
	CreationDate    int
	FirstAlbum      string
	LocationsApi    string `json:"locations"`
	ConcertDatesApi string `json:"concertDates"`
	RelationsApi    string `json:"relations"`
	///////
	Locations    Location
	ConcertDates []Date
	Relations    []Relation
}

type Location struct {
	Id        int
	Locations []string
	Dates     string
}

type Date struct {
	Id    int
	Dates []string
}

type Relation struct {
	Id             int
	DatesLocations map[string][]string
}
