package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"01.gritlab.ax/git/gaddamo/groupie-tracker/functions"
)

var (
	homeTmpl     *template.Template
	artistTmpl   *template.Template
	error400Tmpl *template.Template
	error404Tmpl *template.Template
	error500Tmpl *template.Template
)

func main() {
	var err error
	homeTmpl, err = template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal("ERROR: ", err)
		fmt.Println(homeTmpl)
	}
	artistTmpl, err = template.ParseFiles("static/artist.html")
	if err != nil {
		log.Printf("HTTP %d - ERROR: ARTIST PAGE NOT FOUND: %v", http.StatusInternalServerError, err)

	}
	error400Tmpl, err = template.ParseFiles("static/400.html")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	error404Tmpl, err = template.ParseFiles("static/404.html")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	error500Tmpl, err = template.ParseFiles("static/500.html")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", functions.HomeHandler)
	http.HandleFunc("/artist/", functions.ArtistHandler)

	log.Println("Server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
