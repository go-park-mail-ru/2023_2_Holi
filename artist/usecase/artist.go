package artist_usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
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

func (u *artistUsecase) GetArtistPage(name string) ([]domain.Film, error) {
	films, err := u.atristRepo.GetArtistPage(name)
	if err != nil {
		return nil, err
	}
	logger.Debug("Usecase GetArtistPage:", films)

	return films, nil
}
