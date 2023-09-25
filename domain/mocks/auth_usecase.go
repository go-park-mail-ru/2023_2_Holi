package mocks

import (
	"github.com/stretchr/testify/mock"

	"2023_2_Holi/domain"
)

type AuthUsecase struct {
	mock.Mock
}

func (m *AuthUsecase) Login(credentials domain.Credentials) (domain.Session, error) {
	arguments := m.Called(credentials)

	return arguments.Get(0).(domain.Session), arguments.Error(1)
}

func (m *AuthUsecase) Logout(token string) error {
	arguments := m.Called(token)

	return arguments.Error(0)
}

func (m *AuthUsecase) Register(user domain.User) (int, error) {
	arguments := m.Called(user)

	return arguments.Int(0), arguments.Error(1)
}
