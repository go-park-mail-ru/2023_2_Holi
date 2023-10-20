// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	domain "2023_2_Holi/domain"

	mock "github.com/stretchr/testify/mock"
)

// ArtistUsecase is an autogenerated mock type for the ArtistUsecase type
type ArtistUsecase struct {
	mock.Mock
}

// GetArtistPage provides a mock function with given fields: name, surname
func (_m *ArtistUsecase) GetArtistPage(name string, surname string) ([]domain.Film, error) {
	ret := _m.Called(name, surname)

	var r0 []domain.Film
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) ([]domain.Film, error)); ok {
		return rf(name, surname)
	}
	if rf, ok := ret.Get(0).(func(string, string) []domain.Film); ok {
		r0 = rf(name, surname)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Film)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(name, surname)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewArtistUsecase creates a new instance of ArtistUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewArtistUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *ArtistUsecase {
	mock := &ArtistUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
