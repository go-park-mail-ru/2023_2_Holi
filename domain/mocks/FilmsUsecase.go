// Code generated by mockery v2.35.2. DO NOT EDIT.

package mocks

import (
	domain "2023_2_Holi/domain"

	mock "github.com/stretchr/testify/mock"
)

// FilmsUsecase is an autogenerated mock type for the FilmsUsecase type
type FilmsUsecase struct {
	mock.Mock
}

// GetCastPage provides a mock function with given fields: id
func (_m *FilmsUsecase) GetCastPage(id int) ([]domain.Video, domain.Cast, error) {
	ret := _m.Called(id)

	var r0 []domain.Video
	var r1 domain.Cast
	var r2 error
	if rf, ok := ret.Get(0).(func(int) ([]domain.Video, domain.Cast, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) []domain.Video); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Video)
		}
	}

	if rf, ok := ret.Get(1).(func(int) domain.Cast); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Get(1).(domain.Cast)
	}

	if rf, ok := ret.Get(2).(func(int) error); ok {
		r2 = rf(id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetFilmData provides a mock function with given fields: id
func (_m *FilmsUsecase) GetFilmData(id int) (domain.Video, []domain.Cast, error) {
	ret := _m.Called(id)

	var r0 domain.Video
	var r1 []domain.Cast
	var r2 error
	if rf, ok := ret.Get(0).(func(int) (domain.Video, []domain.Cast, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) domain.Video); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.Video)
	}

	if rf, ok := ret.Get(1).(func(int) []domain.Cast); ok {
		r1 = rf(id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]domain.Cast)
		}
	}

	if rf, ok := ret.Get(2).(func(int) error); ok {
		r2 = rf(id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetFilmsByGenre provides a mock function with given fields: genre
func (_m *FilmsUsecase) GetFilmsByGenre(genre int) ([]domain.Video, error) {
	ret := _m.Called(genre)

	var r0 []domain.Video
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]domain.Video, error)); ok {
		return rf(genre)
	}
	if rf, ok := ret.Get(0).(func(int) []domain.Video); ok {
		r0 = rf(genre)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Video)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(genre)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTopRate provides a mock function with given fields:
func (_m *FilmsUsecase) GetTopRate() (domain.Video, error) {
	ret := _m.Called()

	var r0 domain.Video
	var r1 error
	if rf, ok := ret.Get(0).(func() (domain.Video, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() domain.Video); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.Video)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFilmsUsecase creates a new instance of FilmsUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFilmsUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *FilmsUsecase {
	mock := &FilmsUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
