package domain

type SeriesRepository interface {
	// ListAvailable returns series that have at least 1 episode on disk.
	ListAvailable() ([]*Series, error)
}
