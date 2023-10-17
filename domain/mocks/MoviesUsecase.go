package mocks

import (
	"github.com/stretchr/testify/mock"

	"2023_2_Holi/domain"
)

type MoviesUsecase struct {
	mock.Mock
}

func (m *MoviesUsecase) GetMoviesByGenre(genre string) ([]domain.Movie, error) {
	arguments := m.Called(genre)

	return arguments.Get(0).([]domain.Movie), arguments.Error(1)
}

func (m *MoviesUsecase) GetMovieData(id int) (*domain.Movie, []domain.Artist, error) {
	arguments := m.Called(id)

	return arguments.Get(0).(*domain.Movie), arguments.Get(1).([]domain.Artist), arguments.Error(2)
}
