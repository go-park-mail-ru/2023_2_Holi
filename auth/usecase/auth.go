package usecase

import (
	"time"

	"github.com/google/uuid"

	"2023_2_Holi/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: ur,
	}
}

func (u *userUsecase) Login(user domain.User) (domain.Session, error) {
	expectedUser, err := u.userRepo.GetByName(user.Name)
	if err != nil {
		return domain.Session{}, err
	}

	if expectedUser.Password != user.Password {
		return domain.Session{}, domain.ErrUnauthorized
	}

	return domain.Session{
		Token:       uuid.NewString(),
		SessionData: user.Name,
		ExpiresAt:   time.Now().Add(120 * time.Second),
	}, nil
}
