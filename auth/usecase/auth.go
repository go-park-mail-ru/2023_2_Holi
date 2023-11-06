package auth_usecase

import (
	"bytes"
	"crypto/rand"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

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

func (u *authUsecase) Login(credentials domain.Credentials) (domain.Session, int, error) {
	expectedUser, err := u.authRepo.GetByEmail(credentials.Email)
	if err != nil {
		return domain.Session{}, 0, err
	}
	logs.Logger.Debug("Usecase Login expected user:", expectedUser)

	if !checkPasswords(expectedUser.Password, credentials.Password) {
		return domain.Session{}, 0, domain.ErrWrongCredentials
	}

	session := domain.Session{
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserID:    expectedUser.ID,
	}
	if err = u.sessionRepo.Add(session); err != nil {
		return domain.Session{}, 0, err
	}

	return session, expectedUser.ID, nil
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

	if exists, err := u.authRepo.UserExists(user.Email); exists && err == nil {
		return 0, domain.ErrAlreadyExists
	}

	salt := make([]byte, 8)
	rand.Read(salt)
	user.Password = hashPassword(salt, user.Password)
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

func hashPassword(salt []byte, password []byte) []byte {
	hashedPass := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func checkPasswords(passHash []byte, plainPassword []byte) bool {
	salt := passHash[0:8]
	userPassHash := hashPassword(salt, plainPassword)
	return bytes.Equal(userPassHash, passHash)
}
