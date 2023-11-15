package usecase

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCastPage(t *testing.T) {
	tests := []struct {
		name                     string
		castID                   int
		setFilmsRepoExpectations func(filmsRepo *mocks.FilmsRepository, films []domain.Film, cast domain.Cast)
		good                     bool
	}{
		{
			name:   "GoodCase/Common",
			castID: 1,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, cast domain.Cast) {
				filmsRepo.On("GetCastPage", mock.Anything).Return(films, nil)
				filmsRepo.On("GetCastName", mock.Anything).Return(cast, nil)
			},
			good: true,
		},
		// {
		// 	name:   "ErrorCase/UsecaseError",
		// 	castID: 2,
		// 	setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, cast domain.Cast) {
		// 		filmsRepo.On("GetCastPage", 2).Return(nil, errors.New("Some error"))
		// 		filmsRepo.On("GetCastName", 2).Return(domain.Cast{}, errors.New("Some error"))
		// 	},
		// 	good: false,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fr := new(mocks.FilmsRepository)
			var films []domain.Film
			var cast domain.Cast
			test.setFilmsRepoExpectations(fr, films, cast)

			filmsCase := NewVideoUsecase(fr)
			filmsCaseFilms, filmsCaseCast, err := filmsCase.GetCastPage(test.castID)

			if test.good {
				assert.Nil(t, err)
				assert.EqualValues(t, filmsCaseFilms, films)
				assert.EqualValues(t, filmsCaseCast, cast)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, filmsCaseFilms)
				assert.Empty(t, filmsCaseCast)
			}

			fr.AssertExpectations(t)
		})
	}
}

func TestGetFilmsByGenre(t *testing.T) {
	tests := []struct {
		name                     string
		genre                    string
		setFilmsRepoExpectations func(filmsRepo *mocks.FilmsRepository, films []domain.Film)
		good                     bool
	}{
		{
			name:  "GoodCase/Common",
			genre: "action",
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film) {
				filmsRepo.On("GetFilmsByGenre", "action").Return(films, nil)
			},
			good: true,
		},
		{
			name:  "ErrorCase/UsecaseError",
			genre: "comedy",
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film) {
				filmsRepo.On("GetFilmsByGenre", "comedy").Return(nil, errors.New("Some error"))
			},
			good: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fr := new(mocks.FilmsRepository)
			var films []domain.Film
			test.setFilmsRepoExpectations(fr, films)

			filmsCase := NewVideoUsecase(fr)
			filmsCaseFilms, err := filmsCase.GetFilmsByGenre(test.genre)

			if test.good {
				assert.Nil(t, err)
				assert.EqualValues(t, filmsCaseFilms, films)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, filmsCaseFilms)
			}

			fr.AssertExpectations(t)
		})
	}
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

			filmsCase := NewVideoUsecase(fr)
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
