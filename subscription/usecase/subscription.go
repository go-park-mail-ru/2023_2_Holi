package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"time"
)

type SubsUsecase struct {
	subsRepo domain.SubsRepository
}

func NewSubsUsecase(sr domain.SubsRepository) domain.SubsUsecase {
	return &SubsUsecase{
		subsRepo: sr,
	}
}

func (u *SubsUsecase) Subscribe(subID int) error {
	err := u.subsRepo.Subscribe(subID)
	if err != nil {
		logs.LogError(logs.Logger, "subs_usecase", "Subscribe", err, err.Error())
		return err
	}
	logs.Logger.Debug("Usecase Subscribe:", err)
	return nil
}

func (u *SubsUsecase) UnSubscribe(subID int) error {
	err := u.subsRepo.UnSubscribe(subID)
	if err != nil {
		logs.LogError(logs.Logger, "subs_usecase", "UnSubscribe", err, err.Error())
		return err
	}
	logs.Logger.Debug("Usecase UnSubscribe:", err)
	return nil
}

func (u *SubsUsecase) CheckSub(subID int) (sub time.Time, status bool, error error) {
	timeNow := time.Now()

	subUpTo, err := u.subsRepo.CheckSub(subID)
	if err != nil {
		logs.LogError(logs.Logger, "subs_usecase", "CheckSub", err, err.Error())
		return timeNow, false, err
	}

	logs.Logger.Debug("Usecase CheckSub:", subUpTo)

	if timeNow.Before(subUpTo) || timeNow == subUpTo {
		return subUpTo, true, nil
	} else {
		return timeNow, false, nil
	}
}
