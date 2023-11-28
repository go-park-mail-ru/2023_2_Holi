package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

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
	logs.Logger.Debug("Usecase GetGenres:", genres)

	return genres, nil
}

//func (u *genreUsecase) GetGenresSeries() ([]domain.Genre, error) {
//	genres, err := u.genreRepo.GetGenres()
//	if err != nil {
//		return nil, err
//	}
//	logs.Logger.Debug("Usecase GetGenres:", genres)
//
//	return genres, nil
//}
