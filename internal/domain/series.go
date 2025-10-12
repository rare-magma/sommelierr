package domain

import "time"

type Series struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	Year            int       `json:"year"`
	Overview        string    `json:"overview,omitempty"`
	Poster          string    `json:"poster,omitempty"`
	PosterURL       string    `json:"posterUrl,omitempty"`
	RemotePosterURL string    `json:"remotePosterUrl,omitempty"`
	Images          []Image   `json:"images,omitempty"`
	Added           time.Time `json:"added,omitempty"`
	SourceURL       string
}
