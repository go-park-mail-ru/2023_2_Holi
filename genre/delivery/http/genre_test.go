package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
)

func TestGetGenres(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(usecase *mocks.GenreUsecase)
		status               int
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(usecase *mocks.GenreUsecase) {
				usecase.On("GetGenres").Return([]domain.Genre{}, nil)
			},
			status: http.StatusOK,
		},
		{
			name: "GoodCase/EmptyGenres",
			setUCaseExpectations: func(usecase *mocks.GenreUsecase) {
				usecase.On("GetGenres").Return([]domain.Genre{}, nil)
			},
			status: http.StatusOK,
		},
		{
			name: "ErrorCase/UsecaseError",
			setUCaseExpectations: func(usecase *mocks.GenreUsecase) {
				usecase.On("GetGenres").Return(nil, errors.New("error from usecase"))
			},
			status: http.StatusInternalServerError,
		},
		{
			name: "ErrorCase/InvalidRequest",
			setUCaseExpectations: func(usecase *mocks.GenreUsecase) {
				usecase.On("GetGenres").Return(nil, errors.New("invalid request"))
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := mux.NewRouter()
			mockUsecase := new(mocks.GenreUsecase)
			test.setUCaseExpectations(mockUsecase)

			req, err := http.NewRequest("GET", "/v1/genres", nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			NewGenreHandler(router, mockUsecase)
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}
