package usecase

import (
	"time"

	"github.com/google/uuid"

	"2023_2_Holi/domain"
)

type authUsecase struct {
	authRepo    domain.AuthRepository
	sessionRepo domain.SessionRepository
}

func NewAuthUsecase(ur domain.AuthRepository, sr domain.SessionRepository) domain.AuthUsecase {
	return &authUsecase{
		authRepo:    ur,
		sessionRepo: sr,
	}
}

func (u *authUsecase) Login(user domain.User) (domain.Session, error) {
	expectedUser, err := u.authRepo.GetByName(user.Name)
	if err != nil {
		return domain.Session{}, err
	}

	if expectedUser.Password != user.Password {
		return domain.Session{}, domain.ErrUnauthorized
	}

	session := domain.Session{
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(120 * time.Second),
		UserID:    user.ID,
	}
	if err = u.sessionRepo.Add(session); err != nil {
		return domain.Session{}, nil
	}

	return session, nil
}

func (u *authUsecase) Logout(token string) error {
	if err := u.sessionRepo.DeleteByToken(token); err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) Register(user domain.User) error {
	user.DateJoined = time.Now()
	if err := u.authRepo.AddUser(user); err != nil {
		return err
	}

	return nil
}
