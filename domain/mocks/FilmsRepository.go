// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	domain "2023_2_Holi/domain"

	mock "github.com/stretchr/testify/mock"
)

// FilmsRepository is an autogenerated mock type for the FilmsRepository type
type FilmsRepository struct {
	mock.Mock
}

// GetCastName provides a mock function with given fields: id
func (_m *FilmsRepository) GetCastName(id int) (domain.Cast, error) {
	ret := _m.Called(id)

	var r0 domain.Cast
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (domain.Cast, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) domain.Cast); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.Cast)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCastPage provides a mock function with given fields: id
func (_m *FilmsRepository) GetCastPage(id int) ([]domain.Video, error) {
	ret := _m.Called(id)

	var r0 []domain.Video
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]domain.Video, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) []domain.Video); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Video)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFilmCast provides a mock function with given fields: filmId
func (_m *FilmsRepository) GetFilmCast(filmId int) ([]domain.Cast, error) {
	ret := _m.Called(filmId)

	var r0 []domain.Cast
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]domain.Cast, error)); ok {
		return rf(filmId)
	}
	if rf, ok := ret.Get(0).(func(int) []domain.Cast); ok {
		r0 = rf(filmId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Cast)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(filmId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFilmData provides a mock function with given fields: id
func (_m *FilmsRepository) GetFilmData(id int) (domain.Video, error) {
	ret := _m.Called(id)

	var r0 domain.Video
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (domain.Video, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) domain.Video); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.Video)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFilmsByGenre provides a mock function with given fields: genre
func (_m *FilmsRepository) GetFilmsByGenre(genre string) ([]domain.Video, error) {
	ret := _m.Called(genre)

	var r0 []domain.Video
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.Video, error)); ok {
		return rf(genre)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.Video); ok {
		r0 = rf(genre)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Video)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(genre)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTopRate provides a mock function with given fields:
func (_m *FilmsRepository) GetTopRate() (domain.Video, error) {
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

// NewFilmsRepository creates a new instance of FilmsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFilmsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *FilmsRepository {
	mock := &FilmsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
