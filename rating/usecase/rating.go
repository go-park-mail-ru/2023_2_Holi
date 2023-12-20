package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type RatingUsecase struct {
	RatingRepo domain.RatingRepository
}

func NewRatingUsecase(rr domain.RatingRepository) domain.RatingUsecase {
	return &RatingUsecase{
		RatingRepo: rr,
	}
}

func (u *RatingUsecase) Add(rate domain.Rate) (float64, error) {
	rating, err := u.RatingRepo.Insert(rate)
	if err != nil {
		logs.LogError(logs.Logger, "usecase(rating)", "Add", err, err.Error())
		return 0, err
	}

	return rating, nil
}

func (u *RatingUsecase) Remove(rate domain.Rate) (float64, error) {
	rating, err := u.RatingRepo.Delete(rate)
	if err != nil {
		logs.LogError(logs.Logger, "usecase(rating)", "Remove", err, err.Error())
		return 0, err
	}

	return rating, nil
}

func (u *RatingUsecase) Rated(rate domain.Rate) (bool, int, error) {
	r, rn, err := u.RatingRepo.Exists(rate)
	logs.Logger.Debug("Usecase Rated rated: ", r)
	if err != nil {
		logs.LogError(logs.Logger, "usecase", "Rated", err, err.Error())
		return false, 0, err
	}

	return r, rn, nil
}
