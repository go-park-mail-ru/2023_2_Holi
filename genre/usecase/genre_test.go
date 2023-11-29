package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
)

func TestGetGenres(t *testing.T) {
	tests := []struct {
		name                string
		setRepoExpectations func(repo *mocks.GenreRepository)
		expectedGenres      []domain.Genre
		expectedError       error
	}{
		{
			name: "GoodCase/Common",
			setRepoExpectations: func(repo *mocks.GenreRepository) {
				repo.On("GetGenres").Return([]domain.Genre{
					{
						ID:   1,
						Name: "Action",
					},
					{
						ID:   2,
						Name: "Drama",
					},
				}, nil)
			},
			expectedGenres: []domain.Genre{
				{
					ID:   1,
					Name: "Action",
				},
				{
					ID:   2,
					Name: "Drama",
				},
			},
			expectedError: nil,
		},
		{
			name: "ErrorCase/RepositoryError",
			setRepoExpectations: func(repo *mocks.GenreRepository) {
				repo.On("GetGenres").Return(nil, errors.New("repository error"))
			},
			expectedGenres: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(mocks.GenreRepository)
			test.setRepoExpectations(repo)

			usecase := NewGenreUsecase(repo)

			genres, err := usecase.GetGenres()

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedGenres, genres)
		})
	}
}
