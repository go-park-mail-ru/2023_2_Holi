package usecase

import "2023_2_Holi/domain"

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return userUsecase{
		userRepo: ur,
	}
}
