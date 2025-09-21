package domain

import "time"

type Movie struct {
	Title         string    `json:"title"`
	OriginalTitle string    `json:"originalTitle,omitempty"`
	Year          int       `json:"year"`
	Overview      string    `json:"overview,omitempty"`
	PosterURL     string    `json:"posterUrl,omitempty"`
	Images        []Image   `json:"images,omitempty"`
	Added         time.Time `json:"added,omitempty"`
	SourceURL     string
}