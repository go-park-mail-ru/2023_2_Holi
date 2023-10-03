package collections_usecase

import (
	"2023_2_Holi/domain"
)

type filmUsecase struct {
	filmRepo domain.FilmRepository
}

func NewFilmUsecase(fr domain.FilmRepository) domain.FilmUsecase {
	return &filmUsecase{
		filmRepo: fr,
	}
}

func (u *filmUsecase) GetFilmsByGenre(genre string) ([]domain.Film, error) {
	films, err := u.filmRepo.GetFilmsByGenre(genre)
	if err != nil {
		return nil, err
	}

	return films, nil
}
