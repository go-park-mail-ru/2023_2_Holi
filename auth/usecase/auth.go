package auth_usecase

import (
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type authUsecase struct {
	authRepo    domain.AuthRepository
	sessionRepo domain.SessionRepository
}

func NewAuthUsecase(ar domain.AuthRepository, sr domain.SessionRepository) domain.AuthUsecase {
	return &authUsecase{
		authRepo:    ar,
		sessionRepo: sr,
	}
}

func (u *authUsecase) Login(credentials domain.Credentials) (domain.Session, error) {
	expectedUser, err := u.authRepo.GetByEmail(credentials.Email)
	if err != nil {
		return domain.Session{}, err
	}
	logs.Logger.Debug("Usecase Login expected user:", expectedUser)

	if expectedUser.Password != credentials.Password {
		return domain.Session{}, domain.ErrWrongCredentials
	}

	session := domain.Session{
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserID:    expectedUser.ID,
	}
	if err = u.sessionRepo.Add(session); err != nil {
		return domain.Session{}, err
	}

	return session, nil
}

func (u *authUsecase) Logout(token string) error {
	if token == "" {
		return domain.ErrBadRequest
	}

	if err := u.sessionRepo.DeleteByToken(token); err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) Register(user domain.User) (int, error) {
	if cmp.Equal(user, domain.User{}) {
		return 0, domain.ErrBadRequest
	}

	user.DateJoined = time.Now()

	if exists, err := u.authRepo.UserExists(user.Email); exists && err == nil {
		return 0, domain.ErrAlreadyExists
	}
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

	auth, err := u.sessionRepo.SessionExists(token)
	logs.Logger.Debug("Usecase IsAuth auth: ", auth)
	if err != nil {
		return false, err
	}

	return auth, nil
}
