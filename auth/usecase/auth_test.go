package usecase_test

import (
	"2023_2_Holi/auth/usecase"
	"errors"
	"github.com/bxcodec/faker"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name                       string
		creds                      domain.Credentials
		setAuRepoExpectations      func(creds domain.Credentials, auRepo *mocks.AuthRepository, user *domain.User)
		setSessionRepoExpectations func(sessionRepo *mocks.SessionRepository)
		good                       bool
	}{
		{
			name: "GoodCase/Common",
			creds: domain.Credentials{
				Email:    "uvybini@mail.ru",
				Password: "xrchgvjbk",
			},
			setAuRepoExpectations: func(creds domain.Credentials, auRepo *mocks.AuthRepository, user *domain.User) {
				faker.FakeData(user)
				user.Email = creds.Email
				user.Password = creds.Password
				auRepo.On("GetByEmail", mock.Anything).Return(*user, nil)
			},
			setSessionRepoExpectations: func(sessionRepo *mocks.SessionRepository) {
				sessionRepo.On("Add", mock.Anything).Return(nil)
			},
			good: true,
		},
		{
			name: "BadCase/UserNotFound",
			creds: domain.Credentials{
				Email:    "uvybini@mail.ru",
				Password: "xrchgvjbk",
			},
			setAuRepoExpectations: func(creds domain.Credentials, auRepo *mocks.AuthRepository, user *domain.User) {
				faker.FakeData(user)
				user.Email = creds.Email
				user.Password = creds.Password
				auRepo.On("GetByEmail", mock.Anything).Return(*user, errors.New("some db error"))
			},
			setSessionRepoExpectations: func(sessionRepo *mocks.SessionRepository) {
				sessionRepo.On("Add", mock.Anything).Return(errors.New("another db error")).Maybe()
			},
		},
		{
			name: "BadCase/PasswordDoesntMatch",
			creds: domain.Credentials{
				Email:    "uvybini@mail.ru",
				Password: "xrchgvjbk",
			},
			setAuRepoExpectations: func(creds domain.Credentials, auRepo *mocks.AuthRepository, user *domain.User) {
				faker.FakeData(user)
				user.Email = creds.Email
				user.Password = "ougunorgn"
				auRepo.On("GetByEmail", mock.Anything).Return(*user, nil)
			},
			setSessionRepoExpectations: func(sessionRepo *mocks.SessionRepository) {
				sessionRepo.On("Add", mock.Anything).Return(errors.New("another db error")).Maybe()
			},
		},
		{
			name: "BadCase/InvalidUserId",
			creds: domain.Credentials{
				Email:    "uvybini@mail.ru",
				Password: "xrchgvjbk",
			},
			setAuRepoExpectations: func(creds domain.Credentials, auRepo *mocks.AuthRepository, user *domain.User) {
				faker.FakeData(user)
				user.ID = 0
				user.Email = creds.Email
				user.Password = creds.Password
				auRepo.On("GetByEmail", mock.Anything).Return(*user, nil)
			},
			setSessionRepoExpectations: func(sessionRepo *mocks.SessionRepository) {
				sessionRepo.On("Add", mock.Anything).Return(errors.New("another db error"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			ar := new(mocks.AuthRepository)
			sr := new(mocks.SessionRepository)
			var user domain.User
			test.setAuRepoExpectations(test.creds, ar, &user)
			test.setSessionRepoExpectations(sr)

			auCase := usecase.NewAuthUsecase(ar, sr)
			session, err := auCase.Login(test.creds)

			if test.good {
				assert.Nil(t, err)
				assert.NotEmpty(t, session)
				assert.Equal(t, session.UserID, user.ID)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, session)
			}

			ar.AssertExpectations(t)
			sr.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name                       string
		token                      string
		setSessionRepoExpectations func(sessionRepo *mocks.SessionRepository)
		good                       bool
	}{
		{
			name:  "GoodCase/Common",
			token: uuid.NewString(),
			setSessionRepoExpectations: func(sessionRepo *mocks.SessionRepository) {
				sessionRepo.On("DeleteByToken", mock.Anything).Return(nil)
			},
			good: true,
		},
		{
			name:  "BadCase/EmptyToken",
			token: "",
			setSessionRepoExpectations: func(sessionRepo *mocks.SessionRepository) {
				sessionRepo.On("DeleteByToken", mock.Anything).Return(errors.New("some db error")).Maybe()
			},
			good: false,
		},
		{
			name:  "BadCase/InvalidToken",
			token: "8/refvd 3fdf  sdc",
			setSessionRepoExpectations: func(sessionRepo *mocks.SessionRepository) {
				sessionRepo.On("DeleteByToken", mock.Anything).Return(errors.New("another db error"))
			},
			good: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			ar := new(mocks.AuthRepository)
			sr := new(mocks.SessionRepository)
			test.setSessionRepoExpectations(sr)

			auCase := usecase.NewAuthUsecase(ar, sr)
			err := auCase.Logout(test.token)

			if test.good {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}

			sr.AssertExpectations(t)
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name                       string
		token                      string
		setUserUsecaseExpectations func(u *authUsecase)
		good                       bool
	}{
		{
			name:  "GoodCase/Common",
			token: uuid.NewString(),
			setUserUsecaseExpectations: func(u *authUsecase) {
				sessionRepo.On("DeleteByToken", mock.Anything).Return(nil)
			},
			good: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			ar := new(mocks.AuthRepository)
			sr := new(mocks.SessionRepository)
			test.setSessionRepoExpectations(sr)

			auCase := usecase.NewAuthUsecase(ar, sr)
			err := auCase.Logout(test.token)

			if test.good {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}

			sr.AssertExpectations(t)
		})
	}
}
