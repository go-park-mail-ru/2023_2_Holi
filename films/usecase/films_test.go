package usecase

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetFilmsByGenre(t *testing.T) {
	//TODO: Lexa
}

func TestGetFilmData(t *testing.T) {
	tests := []struct {
		name                     string
		filmID                   int
		setFilmsRepoExpectations func(filmsRepo *mocks.FilmsRepository, film *domain.Film, artists []domain.Cast)
		good                     bool
	}{
		{
			name:   "GoodCase/Common",
			filmID: 1,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, film *domain.Film, artists []domain.Cast) {
				filmsRepo.On("GetFilmData", mock.Anything).Return(*film, nil)
				filmsRepo.On("GetFilmCast", mock.Anything).Return(artists, nil)
			},
			good: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			fr := new(mocks.FilmsRepository)
			var film domain.Film
			var artists []domain.Cast
			test.setFilmsRepoExpectations(fr, &film, artists)

			filmsCase := NewFilmsUsecase(fr)
			filmsCaseFilm, filmsCaseArtists, err := filmsCase.GetFilmData(test.filmID)

			if test.good {
				assert.Nil(t, err)
				assert.EqualValues(t, filmsCaseFilm, film)
				assert.EqualValues(t, filmsCaseArtists, artists)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, filmsCaseFilm)
				assert.Empty(t, filmsCaseArtists)
			}

			fr.AssertExpectations(t)
		})
	}
}
