package http

//
//import (
//	"2023_2_Holi/domain"
//	"github.com/gorilla/context"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//
//	"2023_2_Holi/domain/mocks"
//
//	"github.com/gorilla/mux"
//)
//
//const userID = "1"
//
//func TestAddToFavourites(t *testing.T) {
//	tests := []struct {
//		name            string
//		setExpectations func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase)
//		videoID         string
//		status          int
//	}{
//		{
//			name: "GoodCase/Common",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {
//				uu.On("GetIdBy", mock.Anything).Return(1, nil)
//				fvu.On("AddToFavourites", mock.Anything, mock.Anything).Return(nil)
//			},
//			videoID: "1",
//			status:  http.StatusNoContent,
//		},
//		{
//			name:            "BadCase/InvalidVideoId",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {},
//			videoID:         "ubivgroie",
//			status:          http.StatusBadRequest,
//		},
//		{
//			name:            "BadCase/EmptyVideoId",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {},
//			videoID:         "",
//			status:          http.StatusNotFound,
//		},
//		{
//			name: "BadCase/OutOfRangeVideoId",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {
//				uu.On("GetIdBy", mock.Anything).Return(1, nil)
//				fvu.On("AddToFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: "1234563456789",
//			status:  http.StatusNotFound,
//		},
//		{
//			name: "BadCase/NegativeVideoId",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {
//				uu.On("GetIdBy", mock.Anything).Return(1, nil)
//				fvu.On("AddToFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: "-3",
//			status:  http.StatusNotFound,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			//t.Parallel()
//
//			mfvu := new(mocks.FavouritesUsecase)
//			muu := new(mocks.UtilsUsecase)
//			test.setExpectations(mfvu, muu)
//			req, err := http.NewRequest("POST", "/api/v1/video/favourites/"+test.videoID, nil)
//			assert.NoError(t, err)
//			req.AddCookie(&http.Cookie{
//				Name:     "session_token",
//				Value:    "token",
//				Expires:  time.Now().Add(24 * time.Hour),
//				Path:     "/",
//				HttpOnly: true,
//			})
//			rec := httptest.NewRecorder()
//
//			router := mux.NewRouter()
//			//NewFavouritesHandler(router, mfvu)
//			handler := &FavouritesHandler{
//				FavouritesUsecase: mfvu,
//			}
//			router.HandleFunc("/api/v1/video/favourites/{id}", handler.AddToFavourites).Methods("POST")
//			context.Set(req, "userID", userID)
//			router.ServeHTTP(rec, req)
//
//			assert.Equal(t, test.status, rec.Code)
//			mfvu.AssertExpectations(t)
//			muu.AssertExpectations(t)
//		})
//	}
//}
//
//func TestRemoveFromFavourites(t *testing.T) {
//	tests := []struct {
//		name            string
//		setExpectations func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase)
//		videoID         string
//		status          int
//	}{
//		{
//			name: "GoodCase/Common",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {
//				uu.On("GetIdBy", mock.Anything).Return(1, nil)
//				fvu.On("RemoveFromFavourites", mock.Anything, mock.Anything).Return(nil)
//			},
//			videoID: "1",
//			status:  http.StatusNoContent,
//		},
//		{
//			name:            "BadCase/InvalidVideoId",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {},
//			videoID:         "ubivgroie",
//			status:          http.StatusBadRequest,
//		},
//		{
//			name:            "BadCase/EmptyVideoId",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {},
//			videoID:         "",
//			status:          http.StatusNotFound,
//		},
//		{
//			name: "BadCase/OutOfRangeVideoId",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {
//				uu.On("GetIdBy", mock.Anything).Return(1, nil)
//				fvu.On("RemoveFromFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: "1234563456789",
//			status:  http.StatusNotFound,
//		},
//		{
//			name: "BadCase/NegativeVideoId",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {
//				uu.On("GetIdBy", mock.Anything).Return(1, nil)
//				fvu.On("RemoveFromFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: "-3",
//			status:  http.StatusNotFound,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			//t.Parallel()
//
//			mfvu := new(mocks.FavouritesUsecase)
//			muu := new(mocks.UtilsUsecase)
//			test.setExpectations(mfvu, muu)
//
//			req, err := http.NewRequest("DELETE", "/api/v1/video/favourites/"+test.videoID, nil)
//			assert.NoError(t, err)
//			req.AddCookie(&http.Cookie{
//				Name:     "session_token",
//				Value:    "token",
//				Expires:  time.Now().Add(24 * time.Hour),
//				Path:     "/",
//				HttpOnly: true,
//			})
//			rec := httptest.NewRecorder()
//
//			router := mux.NewRouter()
//			context.Set(req, "userID", userID)
//
//			NewFavouritesHandler(router, mfvu)
//			handler := &FavouritesHandler{
//				FavouritesUsecase: mfvu,
//			}
//			router.HandleFunc("/api/v1/video/favourites/{id}", handler.RemoveFromFavourites).Methods("DELETE")
//			router.ServeHTTP(rec, req)
//
//			assert.Equal(t, test.status, rec.Code)
//			mfvu.AssertExpectations(t)
//			muu.AssertExpectations(t)
//		})
//	}
//}
//
//func TestGetAllFavourites(t *testing.T) {
//	tests := []struct {
//		name            string
//		setExpectations func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase, videos []domain.Video)
//		videos          []domain.Video
//		status          int
//	}{
//		{
//			name: "GoodCase/Common",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase, videos []domain.Video) {
//				uu.On("GetIdBy", mock.Anything).Return(1, nil)
//				fvu.On("GetAllFavourites", mock.Anything).Return(videos, nil)
//			},
//			videos: []domain.Video{
//				domain.Video{
//					ID:               1,
//					Name:             "some",
//					Description:      "desc",
//					PreviewPath:      "path",
//					PreviewVideoPath: "video_path",
//					ReleaseYear:      2007,
//					Rating:           9.5,
//					AgeRestriction:   13,
//				},
//				domain.Video{
//					ID:               2,
//					Name:             "some",
//					Description:      "desc",
//					PreviewPath:      "path",
//					PreviewVideoPath: "video_path",
//					ReleaseYear:      2007,
//					Rating:           9.5,
//					AgeRestriction:   13,
//				},
//			},
//			status: http.StatusOK,
//		},
//		{
//			name: "GoodCase/EmptyFavourites",
//			setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase, videos []domain.Video) {
//				uu.On("GetIdBy", mock.Anything).Return(1, nil)
//				fvu.On("GetAllFavourites", mock.Anything).Return(videos, nil)
//			},
//			videos: []domain.Video{},
//			status: http.StatusOK,
//		},
//		//{
//		//	name:            "BadCase/InvalidVideoId",
//		//	setExpectations: func(fvu *mocks.FavouritesUsecase, uu *mocks.UtilsUsecase) {},
//		//	videoID:         "ubivgroie",
//		//	status:          http.StatusBadRequest,
//		//},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			//t.Parallel()
//
//			mfvu := new(mocks.FavouritesUsecase)
//			muu := new(mocks.UtilsUsecase)
//			test.setExpectations(mfvu, muu, test.videos)
//
//			req, err := http.NewRequest("GET", "/api/v1/video/favourites", nil)
//			assert.NoError(t, err)
//			req.AddCookie(&http.Cookie{
//				Name:     "session_token",
//				Value:    "token",
//				Expires:  time.Now().Add(24 * time.Hour),
//				Path:     "/",
//				HttpOnly: true,
//			})
//			rec := httptest.NewRecorder()
//
//			router := mux.NewRouter()
//			context.Set(req, "userID", userID)
//
//			NewFavouritesHandler(router, mfvu)
//			handler := &FavouritesHandler{
//				FavouritesUsecase: mfvu,
//			}
//			router.HandleFunc("/api/v1/video/favourites", handler.GetAllFavourites).Methods("GET")
//			router.ServeHTTP(rec, req)
//
//			assert.Equal(t, test.status, rec.Code)
//			mfvu.AssertExpectations(t)
//			muu.AssertExpectations(t)
//		})
//	}
//}
