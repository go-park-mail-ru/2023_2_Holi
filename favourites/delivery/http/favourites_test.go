package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"2023_2_Holi/domain/mocks"

	"github.com/gorilla/mux"
)

func TestAddToFavourites(t *testing.T) {
	tests := []struct {
		name            string
		setExpectations func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase)
		videoID         string
		status          int
		good            bool
	}{
		{
			name: "GoodCase/Common",
			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {
				uu.On("GetIdBy", mock.Anything).Return(1, nil)
				fvu.On("Add", mock.Anything).Return(nil)
			},
			videoID: "1",
			status:  http.StatusNoContent,
			good:    true,
		},
		{
			name:            "BadCase/InvalidVideoId",
			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {},
			videoID:         "ubivgroie",
			status:          http.StatusBadRequest,
			good:            true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			router := mux.NewRouter()
			mfvu := new(mocks.FavouritesUsecase)
			muu := new(mocks.UtilsUsecase)
			test.setExpectations(mfvu, muu)

			req, err := http.NewRequest("POST", "/api/v1/video/favourites/"+test.videoID, nil)
			assert.NoError(t, err)
			req.AddCookie(&http.Cookie{
				Name:     "session_token",
				Value:    "token",
				Expires:  time.Now().Add(24 * time.Hour),
				Path:     "/",
				HttpOnly: true,
			})

			rec := httptest.NewRecorder()

			NewFavouritesHandler(router, mfvu, muu)
			handler := &FavouritesHandler{
				FavouritesUsecase: mfvu,
				UtilsUsecase:      muu,
			}
			//handler.AddToFavourites(rec, req)
			//router.ServeHTTP(rec, req)
			router.HandleFunc("/api/v1/video/favourites/{id}", handler.AddToFavourites).Methods("POST")
			router.ServeHTTP(rec, req)
			assert.Equal(t, test.status, rec.Code)
			mfvu.AssertExpectations(t)
			muu.AssertExpectations(t)

		})
	}
}

//func TestGetFilmData(t *testing.T) {
//	tests := []struct {
//		name                 string
//		id                   string
//		setUCaseExpectations func(usecase *mocks.FilmsUsecase, film *domain.Video, artists []domain.Cast, err error)
//		status               int
//	}{
//		{
//			name: "GoodCase/Common",
//			id:   "1",
//			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film *domain.Video, artists []domain.Cast, err error) {
//				usecase.On("GetFilmData", mock.Anything).Return(*film, artists, err)
//			},
//			status: http.StatusOK,
//		},
//		{
//			name: "BadCase/EmptyID",
//			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film *domain.Video, artists []domain.Cast, err error) {
//				usecase.On("GetFilmData", mock.Anything).Return(*film, artists, err)
//			},
//			status: http.StatusNotFound,
//		},
//		{
//			name: "BadCase/WrongIDFormat",
//			id:   "ID",
//			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film *domain.Video, artists []domain.Cast, err error) {
//				usecase.On("GetFilmData", mock.Anything).Return(*film, artists, err)
//			},
//			status: http.StatusBadRequest,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			//t.Parallel()
//
//			router := mux.NewRouter()
//			mockUsecase := new(mocks.FilmsUsecase)
//			var film domain.Video
//			var artists []domain.Cast
//			test.setUCaseExpectations(mockUsecase, &film, artists, nil)
//
//			req, err := http.NewRequest("GET", "/api/v1/films/"+test.id, nil)
//			assert.NoError(t, err)
//
//			rec := httptest.NewRecorder()
//
//			NewVideoHandler(router, mockUsecase)
//
//			handler := &FavouritesHandler{
//				FavouritesUsecase: mockUsecase,
//			}
//
//			router.HandleFunc("/api/v1/films/{id}", handler.GetFilmData).Methods("GET")
//			router.ServeHTTP(rec, req)
//
//			assert.Equal(t, test.status, rec.Code)
//		})
//	}
//}
//
//func TestGetCastPage(t *testing.T) {
//	tests := []struct {
//		name                 string
//		setUCaseExpectations func(usecase *mocks.FilmsUsecase)
//		status               int
//	}{
//		{
//			name: "GoodCase/Common",
//			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
//				usecase.On("GetCastPage", mock.Anything).Return([]domain.Video{}, domain.Cast{}, nil)
//			},
//			status: http.StatusOK,
//		},
//		{
//			name: "GoodCase/EmptyResults",
//			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
//				usecase.On("GetCastPage", mock.Anything).Return([]domain.Video{}, domain.Cast{}, errors.New("error"))
//			},
//			status: http.StatusInternalServerError,
//		},
//		{
//			name: "ErrorCase/UsecaseError",
//			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
//				usecase.On("GetCastPage", mock.Anything).Return([]domain.Video{}, domain.Cast{}, domain.ErrInternalServerError)
//			},
//			status: http.StatusInternalServerError,
//		},
//		{
//			name: "ErrorCase/InvalidRequest",
//			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
//				usecase.On("GetCastPage", mock.Anything).Return([]domain.Video{}, domain.Cast{}, domain.ErrInternalServerError)
//			},
//			status: http.StatusInternalServerError,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			router := mux.NewRouter()
//			mockUsecase := new(mocks.FilmsUsecase)
//			test.setUCaseExpectations(mockUsecase)
//
//			req, err := http.NewRequest("GET", "/api/v1/films/cast/1", nil)
//			assert.NoError(t, err)
//
//			rec := httptest.NewRecorder()
//
//			NewVideoHandler(router, mockUsecase)
//
//			handler := &FavouritesHandler{
//				FavouritesUsecase: mockUsecase,
//			}
//
//			router.HandleFunc("/api/v1/films/cast/{id}", handler.GetCastPage).Methods("GET")
//			router.ServeHTTP(rec, req)
//
//			assert.Equal(t, test.status, rec.Code)
//		})
//	}
//}
