package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type FavouritesUsecase struct {
	favouritesRepo domain.FavouritesRepository
}

func NewFavouritesUsecase(fr domain.FavouritesRepository) domain.FavouritesUsecase {
	return &FavouritesUsecase{
		favouritesRepo: fr,
	}
}

func (u *FavouritesUsecase) Add(videoID, userID int) error {
	err := u.favouritesRepo.Insert(videoID, userID)
	if err != nil {
		logs.LogError(logs.Logger, "usecase", "Add", err, err.Error())
		return err
	}

	return nil
}

func (u *FavouritesUsecase) Remove(videoID, userID int) error {
	err := u.favouritesRepo.Delete(videoID, userID)
	if err != nil {
		logs.LogError(logs.Logger, "usecase", "Remove", err, err.Error())
		return err
	}

	return nil
}

func (u *FavouritesUsecase) GetAll(userID int) ([]domain.Video, error) {
	videos, err := u.favouritesRepo.SelectAll(userID)
	if err != nil {
		logs.LogError(logs.Logger, "usecase", "SelectAll", err, err.Error())
		return []domain.Video{}, err
	}

	return videos, nil
}
