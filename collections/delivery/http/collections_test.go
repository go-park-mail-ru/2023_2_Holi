package collections_http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	collections_http "2023_2_Holi/collections/delivery/http"
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"

	"github.com/gorilla/mux"
)

func TestGetFilmsByGenre(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(usecase *mocks.FilmUsecase)
		status               int
		good                 bool
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(usecase *mocks.FilmUsecase) {
				usecase.On("GetFilmsByGenre", mock.Anything).Return([]domain.Film{}, nil)
			},
			status: http.StatusOK,
			good:   true,
		},
		{
			name: "GoodCase/EmptyFilms",
			setUCaseExpectations: func(usecase *mocks.FilmUsecase) {
				usecase.On("GetFilmsByGenre", mock.Anything).Return([]domain.Film{}, errors.New("error"))
			},
			status: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			router := mux.NewRouter()
			mockUsecase := new(mocks.FilmUsecase)
			test.setUCaseExpectations(mockUsecase)

			req, err := http.NewRequest("GET", "/api/v1/films/genre/{genre}", nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			collections_http.NewFilmHandler(router, mockUsecase)

			handler := &collections_http.FilmHandler{
				FilmUsecase: mockUsecase,
			}

			router.HandleFunc("/api/v1/films/genre/{genre}", handler.GetFilmsByGenre).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}
