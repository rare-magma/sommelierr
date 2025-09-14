package application

import (
	"fmt"
	"math/rand"

	"sommelierr/internal/domain"
)

type GetRandomMovie struct {
	Repo domain.MovieRepository
}
var ErrNoMovies = fmt.Errorf("no available movies found")

func (uc *GetRandomMovie) Execute() (*domain.Movie, error) {
	movies, err := uc.Repo.ListAvailable()
	if err != nil {
		return nil, err
	}
	if len(movies) == 0 {
		return nil, ErrNoMovies
	}
	return movies[rand.Intn(len(movies))], nil
}
