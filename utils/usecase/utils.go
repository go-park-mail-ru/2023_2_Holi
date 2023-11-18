package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type utilsUsecase struct {
	utilsRepo domain.UtilsRepository
}

func NewUtilsUsecase(ur domain.UtilsRepository) domain.UtilsUsecase {
	return &utilsUsecase{
		utilsRepo: ur,
	}
}

func (u *utilsUsecase) GetIdBy(token string) (int, error) {
	id, err := u.utilsRepo.GetIdFromStorage(token)
	if err != nil {
		return 0, err
	}
	logs.Logger.Debug("Utils GetIdBy id:", id)

	return id, nil
}
