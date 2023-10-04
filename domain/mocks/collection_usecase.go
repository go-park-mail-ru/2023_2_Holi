package mocks

import (
	"github.com/stretchr/testify/mock"

	"2023_2_Holi/domain"
)

type FilmUsecase struct {
	mock.Mock
}

func (m *FilmUsecase) GetFilmsByGenre(genre string) ([]domain.Film, error) {
	arguments := m.Called(genre)

	return arguments.Get(0).([]domain.Film), arguments.Error(1)
}
