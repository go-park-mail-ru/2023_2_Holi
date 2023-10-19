package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logs"
)

var logger = logs.LoggerInit()

type artistUsecase struct {
	atristRepo domain.ArtistRepository
}

func NewArtistUsecase(ar domain.ArtistRepository) domain.ArtistUsecase {
	return &artistUsecase{
		atristRepo: ar,
	}
}

func (u *artistUsecase) GetArtistPage(name, surname string) ([]domain.Film, error) {
	films, err := u.atristRepo.GetArtistPage(name, surname)
	if err != nil {
		return nil, err
	}
	logger.Debug("Usecase GetArtistPage:", films)

	return films, nil
}
