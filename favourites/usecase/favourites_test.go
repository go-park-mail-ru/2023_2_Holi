package usecase

//
//import (
//	"2023_2_Holi/domain"
//	"2023_2_Holi/domain/mocks"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"testing"
//)
//
//const userID = 1
//
//func TestAddToFavourites(t *testing.T) {
//	tests := []struct {
//		name            string
//		setExpectations func(fvr *mocks.FavouritesRepository)
//		videoID         int
//		good            bool
//	}{
//		{
//			name: "GoodCase/Common",
//			setExpectations: func(fvr *mocks.FavouritesRepository) {
//				fvr.On("InsertIntoFavourites", mock.Anything, mock.Anything).Return(nil)
//			},
//			videoID: 1,
//			good:    true,
//		},
//		{
//			name: "BadCase/OutOfRangeVideoId",
//			setExpectations: func(fvr *mocks.FavouritesRepository) {
//				fvr.On("InsertIntoFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: 12345634567,
//		},
//		{
//			name: "BadCase/NegativeVideoId",
//			setExpectations: func(fvr *mocks.FavouritesRepository) {
//				fvr.On("InsertIntoFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: -3,
//		},
//		{
//			name: "BadCase/ZeroVideoId",
//			setExpectations: func(fvr *mocks.FavouritesRepository) {
//				fvr.On("InsertIntoFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: 0,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			//t.Parallel()
//
//			fvr := new(mocks.FavouritesRepository)
//			test.setExpectations(fvr)
//
//			fu := NewFavouritesUsecase(fvr)
//			err := fu.AddToFavourites(test.videoID, userID)
//
//			if test.good {
//				assert.Nil(t, err)
//			} else {
//				assert.NotNil(t, err)
//			}
//
//			fvr.AssertExpectations(t)
//		})
//	}
//}
//
//func TestRemoveFromFavourites(t *testing.T) {
//	tests := []struct {
//		name            string
//		setExpectations func(fvr *mocks.FavouritesRepository)
//		videoID         int
//		good            bool
//	}{
//		{
//			name: "GoodCase/Common",
//			setExpectations: func(fvr *mocks.FavouritesRepository) {
//				fvr.On("DeleteFromFavourites", mock.Anything, mock.Anything).Return(nil)
//			},
//			videoID: 1,
//			good:    true,
//		},
//		{
//			name: "BadCase/OutOfRangeVideoId",
//			setExpectations: func(fvr *mocks.FavouritesRepository) {
//				fvr.On("DeleteFromFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: 12345634567,
//		},
//		{
//			name: "BadCase/NegativeVideoId",
//			setExpectations: func(fvr *mocks.FavouritesRepository) {
//				fvr.On("DeleteFromFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: -3,
//		},
//		{
//			name: "BadCase/ZeroVideoId",
//			setExpectations: func(fvr *mocks.FavouritesRepository) {
//				fvr.On("DeleteFromFavourites", mock.Anything, mock.Anything).Return(domain.ErrOutOfRange)
//			},
//			videoID: 0,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			//t.Parallel()
//
//			fvr := new(mocks.FavouritesRepository)
//			test.setExpectations(fvr)
//
//			fu := NewFavouritesUsecase(fvr)
//			err := fu.RemoveFromFavourites(test.videoID, userID)
//
//			if test.good {
//				assert.Nil(t, err)
//			} else {
//				assert.NotNil(t, err)
//			}
//
//			fvr.AssertExpectations(t)
//		})
//	}
//}
//
//func TestGetAllFavourites(t *testing.T) {
//	tests := []struct {
//		name            string
//		setExpectations func(fvr *mocks.FavouritesRepository, vs []domain.Video)
//		videos          []domain.Video
//		good            bool
//	}{
//		{
//			name: "GoodCase/Common",
//			setExpectations: func(fvr *mocks.FavouritesRepository, vs []domain.Video) {
//				fvr.On("SelectAllFavourites", mock.Anything).Return(vs, nil)
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
//					//SeasonsCount:     3,
//				},
//				domain.Video{
//					ID:          2,
//					Name:        "some",
//					Description: "desc",
//					//Duration:         pgtype.Interval{Microseconds: int64(90 * time.Minute), Days: 0, Valid: true},
//					PreviewPath: "path",
//					//MediaPath:        "media_path",
//					PreviewVideoPath: "video_path",
//					ReleaseYear:      2007,
//					Rating:           9.5,
//					AgeRestriction:   13,
//				},
//			},
//			good: true,
//		},
//		{
//			name: "GoodCase/Common",
//			setExpectations: func(fvr *mocks.FavouritesRepository, vs []domain.Video) {
//				fvr.On("SelectAllFavourites", mock.Anything).Return(vs, nil)
//			},
//			videos: []domain.Video{},
//			good:   true,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			//t.Parallel()
//
//			fvr := new(mocks.FavouritesRepository)
//			test.setExpectations(fvr, test.videos)
//
//			fu := NewFavouritesUsecase(fvr)
//			videos, err := fu.GetAllFavourites(userID)
//
//			if test.good {
//				assert.Equal(t, test.videos, videos)
//				assert.Nil(t, err)
//			} else {
//				assert.NotEqual(t, test.videos, videos)
//				assert.NotNil(t, err)
//			}
//
//			fvr.AssertExpectations(t)
//		})
//	}
//}
