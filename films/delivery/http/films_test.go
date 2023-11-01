package films_http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"

	"github.com/gorilla/mux"
)

func TestGetMoviesByGenre(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(usecase *mocks.FilmsUsecase)
		status               int
		good                 bool
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				usecase.On("GetFilmsByGenre", mock.Anything).Return([]domain.Film{}, nil)
			},
			status: http.StatusOK,
			good:   true,
		},
		{
			name: "GoodCase/EmptyFilms",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				usecase.On("GetFilmsByGenre", mock.Anything).Return([]domain.Film{}, errors.New("error"))
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			router := mux.NewRouter()
			mockUsecase := new(mocks.FilmsUsecase)
			test.setUCaseExpectations(mockUsecase)

			req, err := http.NewRequest("GET", "/api/v1/films/genre/{genre}", nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			NewFilmsHandler(router, mockUsecase)

			handler := &FilmsHandler{
				FilmsUsecase: mockUsecase,
			}

			router.HandleFunc("/api/v1/films/genre/{genre}", handler.GetFilmsByGenre).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}

func TestGetFilmData(t *testing.T) {
	tests := []struct {
		name                 string
		id                   string
		setUCaseExpectations func(usecase *mocks.FilmsUsecase, film *domain.Film, artists []domain.Cast, err error)
		status               int
	}{
		{
			name: "GoodCase/Common",
			id:   "1",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film *domain.Film, artists []domain.Cast, err error) {
				usecase.On("GetFilmData", mock.Anything).Return(*film, artists, err)
			},
			status: http.StatusOK,
		},
		{
			name: "BadCase/EmptyID",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film *domain.Film, artists []domain.Cast, err error) {
				usecase.On("GetFilmData", mock.Anything).Return(*film, artists, err)
			},
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/WrongIDFormat",
			id:   "ID",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film *domain.Film, artists []domain.Cast, err error) {
				usecase.On("GetFilmData", mock.Anything).Return(*film, artists, err)
			},
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			router := mux.NewRouter()
			mockUsecase := new(mocks.FilmsUsecase)
			var film domain.Film
			var artists []domain.Cast
			test.setUCaseExpectations(mockUsecase, &film, artists, nil)

			req, err := http.NewRequest("GET", "/api/v1/films/"+test.id, nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			NewFilmsHandler(router, mockUsecase)

			handler := &FilmsHandler{
				FilmsUsecase: mockUsecase,
			}

			router.HandleFunc("/api/v1/films/{id}", handler.GetFilmData).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}

func TestGetCastPage(t *testing.T) {
	tests := []struct {
		name                 string
		id                   string
		setUCaseExpectations func(usecase *mocks.FilmsUsecase, cast *domain.Cast, films []domain.Film, err error)
		status               int
	}{
		{
			name: "GoodCase/Common",
			id:   "1",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, cast *domain.Cast, films []domain.Film, err error) {
				usecase.On("GetCastPage", mock.Anything).Return(films, *cast, err)
			},
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/EmptyID",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, cast *domain.Cast, films []domain.Film, err error) {
				usecase.On("GetCastPage", mock.Anything).Return(films, *cast, err)
			},
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/WrongIDFormat",
			id:   "Wrong id",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, cast *domain.Cast, films []domain.Film, err error) {
				usecase.On("GetCastPage", mock.Anything).Return(films, *cast, err)
			},
			status: http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			router := mux.NewRouter()
			mockUsecase := new(mocks.FilmsUsecase)
			var cast domain.Cast
			var films []domain.Film
			test.setUCaseExpectations(mockUsecase, &cast, films, nil)

			req, err := http.NewRequest("GET", "/api/v1/films/cast/"+test.id, nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			NewFilmsHandler(router, mockUsecase)

			handler := &FilmsHandler{
				FilmsUsecase: mockUsecase,
			}

			router.HandleFunc("/v1/films/cast/{id}", handler.GetFilmData).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}
