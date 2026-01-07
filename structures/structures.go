package structures

type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int     `json:"creationDate"`
	FirstAlbum   string  `json:"firstAlbum"`
	LocationsApi string  `json:"locations"`
	DatesApi     string  `json:"concertDates"`
	RelationApi  string  `json:"relations"`
	events       []Event `json:"-"`
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
	LocationsDates map[string][]string
}
