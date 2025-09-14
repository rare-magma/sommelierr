package domain

type Media struct {
	Title     string `json:"title"`
	Year      int    `json:"year"`
	Overview  string `json:"overview,omitempty"`
	PosterURL string `json:"posterUrl,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
}
