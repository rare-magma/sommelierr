package domain

type Image struct {
	CoverType string `json:"coverType"` // Poster, Banner, Fanart
	URL string `json:"url"` // relative URL
	RemoteURL string `json:"remoteUrl"`
}
