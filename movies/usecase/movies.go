package movies_usecase

import (
	"2023_2_Holi/domain"
)

type moviesUsecase struct {
	movieRepo domain.MoviesRepository
}

func NewMoviesUsecase(fr domain.MoviesRepository) domain.MoviesUsecase {
	return &moviesUsecase{
		movieRepo: fr,
	}
}

func (u *moviesUsecase) GetMoviesByGenre(genre string) ([]domain.Movie, error) {
	films, err := u.movieRepo.GetMoviesByGenre(genre)
	if err != nil {
		return nil, err
	}

	return films, nil
}
