// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/coutcout/covoiturage-csvreader/domain"

	mock "github.com/stretchr/testify/mock"
)

// JourneyRepositoryInterface is an autogenerated mock type for the JourneyRepositoryInterface type
type JourneyRepositoryInterface struct {
	mock.Mock
}

// Add provides a mock function with given fields: journey
func (_m *JourneyRepositoryInterface) Add(journey *domain.Journey) (bool, error) {
	ret := _m.Called(journey)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*domain.Journey) (bool, error)); ok {
		return rf(journey)
	}
	if rf, ok := ret.Get(0).(func(*domain.Journey) bool); ok {
		r0 = rf(journey)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*domain.Journey) error); ok {
		r1 = rf(journey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewJourneyRepositoryInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewJourneyRepositoryInterface creates a new instance of JourneyRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJourneyRepositoryInterface(t mockConstructorTestingTNewJourneyRepositoryInterface) *JourneyRepositoryInterface {
	mock := &JourneyRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
