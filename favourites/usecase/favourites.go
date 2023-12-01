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

func (u *FavouritesUsecase) AddToFavourites(videoID, userID int) error {
	err := u.favouritesRepo.InsertIntoFavourites(videoID, userID)
	if err != nil {
		logs.LogError(logs.Logger, "usecase", "AddToFavourites", err, err.Error())
		return err
	}

	return nil
}

func (u *FavouritesUsecase) RemoveFromFavourites(videoID, userID int) error {
	err := u.favouritesRepo.DeleteFromFavourites(videoID, userID)
	if err != nil {
		logs.LogError(logs.Logger, "usecase", "RemoveFromFavourites", err, err.Error())
		return err
	}

	return nil
}

func (u *FavouritesUsecase) GetAllFavourites(userID int) ([]domain.Video, error) {
	videos, err := u.favouritesRepo.SelectAllFavourites(userID)
	if err != nil {
		logs.LogError(logs.Logger, "usecase", "GetAllFavourites", err, err.Error())
		return []domain.Video{}, err
	}

	return videos, nil
}

func (u *FavouritesUsecase) Favourite(videoID, userID int) (bool, error) {
	f, err := u.favouritesRepo.Exists(videoID, userID)
	logs.Logger.Debug("Usecase Favourite favourite: ", f)
	if err != nil {
		logs.LogError(logs.Logger, "usecase", "GetAllFFavouriteavourites", err, err.Error())
		return false, err
	}

	return f, nil
}
