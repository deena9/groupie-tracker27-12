package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func LocationCount(locations []string) int {
	LocationList := make(map[string]string)
	for _, loc := range locations {
		split := strings.Split(loc, "-")
		if len(split) == 2 {
			LocationList[split[1]] = split[0]
		}
	}
	return len(LocationList)
}

func FetchArtists() error {
	resp, err := http.Get(artistsEndpoint)
	if err != nil {
		return fmt.Errorf("HTTP %d - failed to fetch artists: %v", http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &artists); err != nil {
		return err
	}
	log.Println("Artists fetched successfully")
	return nil
}

func FetchLocations() error {
	resp, err := http.Get(locationsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiResponse LocationAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}
	locations = apiResponse.Index
	log.Println("Locations fetched successfully")
	return nil
}

func FetchDates() error {
	resp, err := http.Get(datesEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Use the wrapper struct to unmarshal the data
	var apiResponse DatesAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}

	// Assign the fetched dates
	dates = apiResponse.Index
	log.Println("Dates fetched successfully")
	return nil
}

func FetchRelations() error {
	resp, err := http.Get(relationsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &relations); err != nil {
		return err
	}
	log.Println("Relations fetched successfully")
	return nil
}

// Fetch All Data at Once
func FetchAllData() {
	if len(artists) == 0 {
		if err := FetchArtists(); err != nil {
			log.Fatal("Error fetching artists:", err)
		}
		if err := FetchLocations(); err != nil {
			log.Fatal("Error fetching locations:", err)
		}
		if err := FetchDates(); err != nil {
			log.Fatal("Error fetching dates:", err)
		}
		if err := FetchRelations(); err != nil {
			log.Fatal("Error fetching relations:", err)
		}
	}
}

func FetchArtistData(id int) (Artist, []string, []string, map[string][]string, error) {
	// Fetch the artist by ID
	var selectedArtist Artist
	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			break
		}
	}

	// Fetch associated data (locations, dates, and relations)
	var associatedLocations []string
	var associatedDates []string
	var associatedRelations map[string][]string
	var locationCount int

	for _, location := range locations {
		if location.ID == id {
			for _, loc := range location.Locations {
				cleanLoc := strings.ReplaceAll(loc, "_", " ")
				cleanLoc = strings.ReplaceAll(cleanLoc, "-", " ")
				associatedLocations = append(associatedLocations, cleanLoc)
			}
			locationCount = LocationCount(associatedLocations)
		}
	}
	selectedArtist.LocationCount = locationCount

	for _, date := range dates {
		if date.ID == id {
			for _, d := range date.Dates {
				cleanDate := strings.ReplaceAll(d, "*", "")
				associatedDates = append(associatedDates, cleanDate)
			}
		}
	}

	for _, relation := range relations.Index {
		if relation.ID == id {
			associatedRelations = make(map[string][]string)
			for loc, relDates := range relation.DatesLocations {
				cleanLoc := strings.ReplaceAll(loc, "_", " ")
				cleanLoc = strings.ReplaceAll(cleanLoc, "-", " ")
				cleanRelDates := []string{}
				for _, d := range relDates {
					cleanDate := strings.ReplaceAll(d, "*", "")
					cleanRelDates = append(cleanRelDates, cleanDate)
				}
				associatedRelations[cleanLoc] = cleanRelDates
			}
		}
	}

	return selectedArtist, associatedLocations, associatedDates, associatedRelations, nil
}
