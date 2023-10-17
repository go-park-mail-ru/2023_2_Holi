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

func (u *moviesUsecase) GetMovieData(id int) (*domain.Movie, []domain.Artist, error) {
	movie, err := u.movieRepo.GetMovieData(id)
	if err != nil {
		return nil, nil, err
	}
	artists, err := u.movieRepo.GetMovieArtists(id)
	if err != nil {
		return nil, nil, err
	}

	return movie, artists, nil
}
