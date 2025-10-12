package domain

import "time"

type Movie struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	OriginalTitle   string    `json:"originalTitle,omitempty"`
	Year            int       `json:"year"`
	Overview        string    `json:"overview,omitempty"`
	Poster          string    `json:"poster,omitempty"`
	PosterURL       string    `json:"posterUrl,omitempty"`
	RemotePosterURL string    `json:"remotePosterUrl,omitempty"`
	Images          []Image   `json:"images,omitempty"`
	Added           time.Time `json:"added,omitempty"`
	SourceURL       string
}
