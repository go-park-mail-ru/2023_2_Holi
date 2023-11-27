// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	domain "2023_2_Holi/domain"

	mock "github.com/stretchr/testify/mock"
)

// SearchRepository is an autogenerated mock type for the SearchRepository type
type SearchRepository struct {
	mock.Mock
}

// GetSuitableCast provides a mock function with given fields: searchStr
func (_m *SearchRepository) GetSuitableCast(searchStr string) ([]domain.Cast, error) {
	ret := _m.Called(searchStr)

	var r0 []domain.Cast
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.Cast, error)); ok {
		return rf(searchStr)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.Cast); ok {
		r0 = rf(searchStr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Cast)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(searchStr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSuitableFilms provides a mock function with given fields: searchStr
func (_m *SearchRepository) GetSuitableFilms(searchStr string) ([]domain.Video, error) {
	ret := _m.Called(searchStr)

	var r0 []domain.Video
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.Video, error)); ok {
		return rf(searchStr)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.Video); ok {
		r0 = rf(searchStr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Video)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(searchStr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSearchRepository creates a new instance of SearchRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSearchRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *SearchRepository {
	mock := &SearchRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}