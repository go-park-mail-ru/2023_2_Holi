package movies_http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	movies_http "2023_2_Holi/movies/delivery/http"

	"github.com/gorilla/mux"
)

func TestGetMoviesByGenre(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(usecase *mocks.MoviesUsecase)
		status               int
		good                 bool
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(usecase *mocks.MoviesUsecase) {
				usecase.On("GetFilmsByGenre", mock.Anything).Return([]domain.Movie{}, nil)
			},
			status: http.StatusOK,
			good:   true,
		},
		{
			name: "GoodCase/EmptyFilms",
			setUCaseExpectations: func(usecase *mocks.MoviesUsecase) {
				usecase.On("GetFilmsByGenre", mock.Anything).Return([]domain.Movie{}, errors.New("error"))
			},
			status: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			router := mux.NewRouter()
			mockUsecase := new(mocks.MoviesUsecase)
			test.setUCaseExpectations(mockUsecase)

			req, err := http.NewRequest("GET", "/api/v1/films/genre/{genre}", nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			movies_http.NewMoviesHandler(router, mockUsecase)

			handler := &movies_http.MoviesHandler{
				MoviesUsecase: mockUsecase,
			}

			router.HandleFunc("/api/v1/films/genre/{genre}", handler.GetMoviesByGenre).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}
