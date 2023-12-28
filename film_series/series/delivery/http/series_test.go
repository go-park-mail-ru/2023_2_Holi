package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"

	"github.com/gorilla/mux"
)

func TestGetSeriesByGenre(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(usecase *mocks.SeriesUsecase)
		status               int
		good                 bool
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase) {
				usecase.On("GetSeriesByGenre", mock.Anything).Return([]domain.Video{}, nil)
			},
			status: http.StatusOK,
			good:   true,
		},
		{
			name: "GoodCase/EmptySeries",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase) {
				usecase.On("GetSeriesByGenre", mock.Anything).Return([]domain.Video{}, errors.New("error"))
			},
			status: http.StatusInternalServerError,
		},
		{
			name: "GoodCase/NonEmptySeries",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase) {
				usecase.On("GetSeriesByGenre", mock.Anything).Return([]domain.Video{{ID: 1, Name: "Video 1"}, {ID: 2, Name: "Video 2"}}, nil)
			},
			status: http.StatusOK,
			good:   true,
		},
		{
			name: "ErrorCase/UsecaseError",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase) {
				usecase.On("GetSeriesByGenre", mock.Anything).Return(nil, errors.New("error from usecase"))
			},
			status: http.StatusInternalServerError,
		},
		{
			name: "ErrorCase/InvalidRequest",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase) {
				usecase.On("GetSeriesByGenre", mock.Anything).Return(nil, errors.New("invalid request"))
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := mux.NewRouter()
			mockUsecase := new(mocks.SeriesUsecase)
			test.setUCaseExpectations(mockUsecase)

			genreID := 1

			req, err := http.NewRequest("GET", "/api/v1/series/genre/"+strconv.Itoa(genreID), nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			NewSeriesHandler(router, mockUsecase)

			handler := &SeriesHandler{
				SeriesUsecase: mockUsecase,
			}

			router.HandleFunc("/api/v1/series/genre/{id}", handler.GetSeriesByGenre).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}

func TestGetSeriesData(t *testing.T) {
	tests := []struct {
		name                 string
		id                   string
		setUCaseExpectations func(usecase *mocks.SeriesUsecase, series *domain.Video, artists []domain.Cast, episodes []domain.Episode, err error)
		status               int
	}{
		{
			name: "GoodCase/Common",
			id:   "1",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase, series *domain.Video, artists []domain.Cast, episodes []domain.Episode, err error) {
				usecase.On("GetSeriesData", mock.Anything).Return(*series, artists, episodes, err)
			},
			status: http.StatusOK,
		},
		{
			name: "BadCase/EmptyID",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase, series *domain.Video, artists []domain.Cast, episodes []domain.Episode, err error) {
				usecase.On("GetSeriesData", mock.Anything).Return(*series, artists, episodes, err)
			},
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/WrongIDFormat",
			id:   "ID",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase, series *domain.Video, artists []domain.Cast, episodes []domain.Episode, err error) {
				usecase.On("GetSeriesData", mock.Anything).Return(*series, artists, episodes, err)
			},
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := mux.NewRouter()
			mockUsecase := new(mocks.SeriesUsecase)
			var series domain.Video
			var artists []domain.Cast
			var episodes []domain.Episode
			test.setUCaseExpectations(mockUsecase, &series, artists, episodes, nil)

			req, err := http.NewRequest("GET", "/api/v1/series/"+test.id, nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			NewSeriesHandler(router, mockUsecase)

			handler := &SeriesHandler{
				SeriesUsecase: mockUsecase,
			}

			router.HandleFunc("/api/v1/series/{id}", handler.GetSeriesData).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}

func TestGetCastPageSeries(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(usecase *mocks.SeriesUsecase)
		status               int
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase) {
				usecase.On("GetCastPageSeries", mock.Anything).Return([]domain.Video{}, domain.Cast{}, nil)
			},
			status: http.StatusOK,
		},
		{
			name: "GoodCase/EmptyResults",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase) {
				usecase.On("GetCastPageSeries", mock.Anything).Return([]domain.Video{}, domain.Cast{}, errors.New("error"))
			},
			status: http.StatusInternalServerError,
		},
		{
			name: "ErrorCase/UsecaseError",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase) {
				usecase.On("GetCastPageSeries", mock.Anything).Return([]domain.Video{}, domain.Cast{}, errors.New("error from usecase"))
			},
			status: http.StatusInternalServerError,
		},
		{
			name: "ErrorCase/InvalidRequest",
			setUCaseExpectations: func(usecase *mocks.SeriesUsecase) {
				usecase.On("GetCastPageSeries", mock.Anything).Return([]domain.Video{}, domain.Cast{}, errors.New("invalid request"))
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := mux.NewRouter()
			mockUsecase := new(mocks.SeriesUsecase)
			test.setUCaseExpectations(mockUsecase)

			req, err := http.NewRequest("GET", "/api/v1/series/cast/1", nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			NewSeriesHandler(router, mockUsecase)

			handler := &SeriesHandler{
				SeriesUsecase: mockUsecase,
			}

			router.HandleFunc("/api/v1/series/cast/{id}", handler.GetCastPageSeries).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}
