package films_usecase_test

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	films_usecase "2023_2_Holi/films/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//
// "errors"
// "testing"

// "github.com/bxcodec/faker"
// "github.com/google/uuid"
// "github.com/stretchr/testify/assert"
// "github.com/stretchr/testify/mock"

// "2023_2_Holi/domain"
// "2023_2_Holi/domain/mocks"

func TestGetFilmsByGenre(t *testing.T) {
	//TODO: Lexa
}

func TestGetFilmData(t *testing.T) {
	tests := []struct {
		name                     string
		filmID                   int
		setFilmsRepoExpectations func(filmsRepo *mocks.FilmsRepository, film *domain.Film, artists []domain.Artist)
		good                     bool
	}{
		{
			name:   "GoodCase/Common",
			filmID: 1,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, film *domain.Film, artists []domain.Artist) {
				filmsRepo.On("GetFilmData", mock.Anything).Return(*film, nil)
				filmsRepo.On("GetFilmArtists", mock.Anything).Return(artists, nil)
			},
			good: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			fr := new(mocks.FilmsRepository)
			var film domain.Film
			var artists []domain.Artist
			test.setFilmsRepoExpectations(fr, &film, artists)

			filmsCase := films_usecase.NewFilmsUsecase(fr)
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
