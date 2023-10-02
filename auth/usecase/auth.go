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

func (u *authUsecase) Login(credentials domain.Credentials) (domain.Session, error) {
	expectedUser, err := u.authRepo.GetByName(credentials.Name)
	if err != nil {
		return domain.Session{}, domain.ErrNotFound
	}

	if expectedUser.Password != credentials.Password {
		return domain.Session{}, domain.ErrWrongCredentials
	}

	session := domain.Session{
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserID:    expectedUser.ID,
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

func (u *authUsecase) Register(user domain.User) (int, error) {
	user.DateJoined = time.Now()

	if id, err := u.authRepo.AddUser(user); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func (u *authUsecase) IsAuth(token string) (bool, error) {
	if token == "" {
		return false, domain.ErrBadRequest
	}

	auth, err := u.sessionRepo.IsAuth(token)
	if err != nil {
		return false, err
	}

	return auth, nil
}
