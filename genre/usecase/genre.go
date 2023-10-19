package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logs"
)

var logger = logs.LoggerInit()

type genreUsecase struct {
	genreRepo domain.GenreRepository
}

func NewGenreUsecase(gr domain.GenreRepository) domain.GenreUsecase {
	return &genreUsecase{
		genreRepo: gr,
	}
}

func (u *genreUsecase) GetGenres() ([]domain.Genre, error) {
	genres, err := u.genreRepo.GetGenres()
	if err != nil {
		return nil, err
	}
	logger.Debug("Usecase GetGenres:", genres)

	return genres, nil
}
