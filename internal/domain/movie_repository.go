package domain

type MovieRepository interface {
	// ListAvailable returns movies that have a file on disk.
	ListAvailable() ([]*Movie, error)
	GetPoster(int, string) (string, error)
}
