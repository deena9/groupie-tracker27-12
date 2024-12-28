package functions

import (
	"text/template"
)

var (
	tpl        *template.Template
	fetchError bool
)

const (
	apiURL            = "https://groupietrackers.herokuapp.com/api"
	artistsEndpoint   = apiURL + "/artists"
	locationsEndpoint = apiURL + "/locations"
	datesEndpoint     = apiURL + "/dates"
	relationsEndpoint = apiURL + "/relation"
)

// Artist data structure remains the same
type Artist struct {
	ID            int      `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Members       []string `json:"members"`
	CreationDate  int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	LocationCount int
}

// Artist data structure remains the same
type ArtistPageData struct {
	ID            int      `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Members       []string `json:"members"`
	CreationDate  int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	Locations     []string
	Dates         []string
	Relations     []string
	LocationCount int
}

type LocationAPIResponse struct {
	Index []Location `json:"index"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type DatesAPIResponse struct {
	Index []Dates `json:"index"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

/* type Relation struct {
	Index []struct {
		ID        int      `json:"id"`
		Dates     []string `json:"dates"`
		Locations []string `json:"locations"`
	} `json:"index"`
} */

type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type ErrorResponse struct {
	Code    int
	Message string
}

var (
	artists   []Artist
	locations []Location
	dates     []Dates
	relations Relation
)
