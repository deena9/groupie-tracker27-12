package functions

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// ...existing code...
var artistTmpl = template.Must(template.ParseFiles("static/artist.html"))
var homeTmpl = template.Must(template.ParseFiles("static/index.html"))
var error400Tmpl = template.Must(template.ParseFiles("static/400.html"))
var error404Tmpl = template.Must(template.ParseFiles("static/404.html"))
var error500Tmpl = template.Must(template.ParseFiles("static/500.html"))

// ...existing code...

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	FetchAllData()

	if fetchError {
		HandleError(w, http.StatusInternalServerError, "Failed to fetch data")
		return
	}

	artistID := r.URL.Path[len("/artist/"):]
	if artistID == "" {
		HandleError(w, http.StatusBadRequest, "Artist ID is required")
		return
	}

	idString := ""
	if _, err := fmt.Sscanf(artistID, "%s", &idString); err != nil {
		HandleError(w, http.StatusBadRequest, "Invalid artist ID format")
		return
	}

	id, iderr := strconv.Atoi(idString)
	if iderr != nil {
		HandleError(w, http.StatusBadRequest, "Invalid artist ID format")
		return
	}

	// Check if artist exists
	artistFound := false
	for _, artist := range artists {
		if artist.ID == id {
			artistFound = true
			break
		}
	}

	if !artistFound {
		HandleError(w, http.StatusNotFound, "Artist not found")
		return
	}

	// Fetch the artist and related data
	artist, locations, dates, relations, err := FetchArtistData(id)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Failed to fetch artist data")
		return
	}

	stringRelations := []string{}
	for location, dates := range relations {
		for _, date := range dates {
			stringRelations = append(stringRelations, location+" "+date)
		}
	}
	locationCount := 0

	APD := ArtistPageData{
		ID:            artist.ID,
		Image:         artist.Image,
		Name:          artist.Name,
		Members:       artist.Members,
		CreationDate:  artist.CreationDate,
		FirstAlbum:    artist.FirstAlbum,
		Locations:     locations,
		Dates:         dates,
		Relations:     stringRelations,
		LocationCount: locationCount,
	}

	if err := artistTmpl.Execute(w, APD); err != nil {
		HandleError(w, http.StatusInternalServerError, "Failed to render template")
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	// switch r.Method {
	// case "GET":
	// 	// Serve the initial form on GET request
	// 	err := tpl.ExecuteTemplate(w, "index.html", nil)
	// 	if err != nil {
	// 		HandleError(w, http.StatusInternalServerError, "500.html")
	// 	}

	FetchAllData()
	if fetchError {
		HandleError(w, http.StatusInternalServerError, "Failed to fetch data")
		return
	}

	if err := homeTmpl.Execute(w, artists); err != nil {
		HandleError(w, http.StatusInternalServerError, "Failed to render template")
		return
	}
}

// }

func HandleError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	log.Printf("Error %d: %s", code, message) // Log the error

	switch code {
	case http.StatusBadRequest: // 400
		if err := error400Tmpl.Execute(w, ErrorResponse{Code: code, Message: message}); err != nil {
			http.Error(w, message, code)
		}
	case http.StatusNotFound: // 404
		if err := error404Tmpl.Execute(w, ErrorResponse{Code: code, Message: message}); err != nil {
			http.Error(w, message, code)
		}
	case http.StatusInternalServerError: // 500
		if err := error500Tmpl.Execute(w, ErrorResponse{Code: code, Message: message}); err != nil {
			http.Error(w, message, code)
		}
	default:
		http.Error(w, message, code)
	}
}
