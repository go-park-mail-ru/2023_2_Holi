package mocks

import (
	"github.com/stretchr/testify/mock"

	"2023_2_Holi/domain"
)

type FilmsUsecase struct {
	mock.Mock
}

func (m *FilmsUsecase) GetFilmsByGenre(genre string) ([]domain.Film, error) {
	arguments := m.Called(genre)

	return arguments.Get(0).([]domain.Film), arguments.Error(1)
}

func (m *FilmsUsecase) GetFilmData(id int) (*domain.Film, []domain.Cast, error) {
	arguments := m.Called(id)

	return arguments.Get(0).(*domain.Film), arguments.Get(1).([]domain.Cast), arguments.Error(2)
}
