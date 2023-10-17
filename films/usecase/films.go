package films_usecase

import (
	"2023_2_Holi/domain"
)

type filmsUsecase struct {
	filmRepo domain.FilmsRepository
}

func NewFilmsUsecase(fr domain.FilmsRepository) domain.FilmsUsecase {
	return &filmsUsecase{
		filmRepo: fr,
	}
}

func (u *filmsUsecase) GetFilmsByGenre(genre string) ([]domain.Film, error) {
	films, err := u.filmRepo.GetFilmsByGenre(genre)
	if err != nil {
		return nil, err
	}

	return films, nil
}

func (u *filmsUsecase) GetFilmData(id int) (*domain.Film, []domain.Artist, error) {
	film, err := u.filmRepo.GetFilmData(id)
	if err != nil {
		return nil, nil, err
	}
	artists, err := u.filmRepo.GetFilmArtists(id)
	if err != nil {
		return nil, nil, err
	}

	return film, artists, nil
}
