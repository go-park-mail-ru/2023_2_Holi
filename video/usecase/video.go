package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type videoUsecase struct {
	videoRepo domain.VideoRepository
}

func NewVideoUsecase(vr domain.VideoRepository) domain.VideoUsecase {
	return &videoUsecase{
		videoRepo: vr,
	}
}

func (u *videoUsecase) AddToFavourites(id int) error {
	err := u.videoRepo.AddToFavourites(id)
	if err != nil {
		logs.LogError(logs.Logger, "usecase", "AddToFavourites", err, err.Error())
		return domain.Film{}, nil, err
	}

	return nil
}
