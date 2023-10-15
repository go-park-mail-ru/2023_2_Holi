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
