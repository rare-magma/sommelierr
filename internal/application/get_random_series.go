package application

import (
	"fmt"
	"math/rand"

	"sommelierr/internal/domain"
)

type GetRandomSeries struct {
	Repo domain.SeriesRepository
}

var ErrNoSeries = fmt.Errorf("no available series found")

func (uc *GetRandomSeries) Execute() (*domain.Series, error) {
	series, err := uc.Repo.ListAvailable()
	if err != nil {
		return nil, err
	}
	if len(series) == 0 {
		return nil, ErrNoSeries
	}
	randomPick := series[rand.Intn(len(series))]
	if randomPick.PosterURL != "" {
		p, err := uc.Repo.GetPoster(randomPick.Id, randomPick.PosterURL)
		if err != nil {
			return nil, err
		}
		randomPick.Poster = p
	}
	return randomPick, nil
}
