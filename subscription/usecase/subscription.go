package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type subsUsecase struct {
	subsRepo domain.SubsRepository
}

func NewSubsUsecase(sr domain.SubsRepository) domain.SubsUsecase {
	return &subsUsecase{
		subsRepo: sr,
	}
}

func (u *subsUsecase) Subscribe(subID int, flag int) error {
	err := u.subsRepo.Subscribe(subID, flag)
	if err != nil {
		logs.LogError(logs.Logger, "subs_usecase", "Subscribe", err, err.Error())
		return err
	}
	logs.Logger.Debug("Usecase Subscribe:", err)
	return nil
}

func (u *subsUsecase) UnSubscribe(subID int) error {
	err := u.subsRepo.UnSubscribe(subID)
	if err != nil {
		logs.LogError(logs.Logger, "subs_usecase", "UnSubscribe", err, err.Error())
		return err
	}
	logs.Logger.Debug("Usecase UnSubscribe:", err)
	return nil
}

func (u *subsUsecase) CheckSub(subID int) error {
	err := u.subsRepo.CheckSub(subID)
	if err != nil {
		logs.LogError(logs.Logger, "subs_usecase", "CheckSub", err, err.Error())
		return err
	}
	logs.Logger.Debug("Usecase CheckSub:", err)
	return nil
}

func (u *subsUsecase) GetSubInfo(subID int) (domain.SubInfo, error) {
	subInfo, err := u.subsRepo.GetSubInfo(subID)
	if err != nil {
		logs.LogError(logs.Logger, "subs_usecase", "GetSubInfo", err, err.Error())
		return domain.SubInfo{}, err
	}
	logs.Logger.Debug("Usecase GetSubInfo:", err)
	return subInfo, nil
}
