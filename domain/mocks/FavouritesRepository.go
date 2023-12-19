// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	domain "2023_2_Holi/domain"

	mock "github.com/stretchr/testify/mock"
)

// FavouritesRepository is an autogenerated mock type for the FavouritesRepository type
type FavouritesRepository struct {
	mock.Mock
}

// DeleteFromFavourites provides a mock function with given fields: videoID, userID
func (_m *FavouritesRepository) DeleteFromFavourites(videoID int, userID int) error {
	ret := _m.Called(videoID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, int) error); ok {
		r0 = rf(videoID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: videoID, userID
func (_m *FavouritesRepository) Exists(videoID int, userID int) (bool, error) {
	ret := _m.Called(videoID, userID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) (bool, error)); ok {
		return rf(videoID, userID)
	}
	if rf, ok := ret.Get(0).(func(int, int) bool); ok {
		r0 = rf(videoID, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(videoID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertIntoFavourites provides a mock function with given fields: videoID, userID
func (_m *FavouritesRepository) InsertIntoFavourites(videoID int, userID int) error {
	ret := _m.Called(videoID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, int) error); ok {
		r0 = rf(videoID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SelectAllFavourites provides a mock function with given fields: userID
func (_m *FavouritesRepository) SelectAllFavourites(userID int) ([]domain.Video, error) {
	ret := _m.Called(userID)

	var r0 []domain.Video
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]domain.Video, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(int) []domain.Video); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Video)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFavouritesRepository creates a new instance of FavouritesRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFavouritesRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *FavouritesRepository {
	mock := &FavouritesRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
