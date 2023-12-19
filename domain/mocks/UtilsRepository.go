// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UtilsRepository is an autogenerated mock type for the UtilsRepository type
type UtilsRepository struct {
	mock.Mock
}

// GetIdFromStorage provides a mock function with given fields: token
func (_m *UtilsRepository) GetIdFromStorage(token string) (int, error) {
	ret := _m.Called(token)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUtilsRepository creates a new instance of UtilsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUtilsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UtilsRepository {
	mock := &UtilsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
