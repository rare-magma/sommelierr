package domain

import "time"

type Movie struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	OriginalTitle string    `json:"originalTitle,omitempty"`
	Year          int       `json:"year"`
	Overview      string    `json:"overview,omitempty"`
	PosterURL     string    `json:"posterUrl,omitempty"`
	Images        []Image   `json:"images,omitempty"`
	Added         time.Time `json:"added,omitempty"`
}

type MovieRepository interface {
	// ListAvailable returns movies that have a file on disk.
	ListAvailable() ([]*Movie, error)
}
