// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	io "io"

	domain "github.com/coutcout/covoiturage-csvreader/domain"

	mock "github.com/stretchr/testify/mock"
)

// JourneyParser is an autogenerated mock type for the JourneyParser type
type JourneyParser struct {
	mock.Mock
}

// Parse provides a mock function with given fields: reader, journeyChan, errorChan
func (_m *JourneyParser) Parse(reader io.Reader, journeyChan chan<- *domain.Journey, errorChan chan<- string) {
	_m.Called(reader, journeyChan, errorChan)
}

type mockConstructorTestingTNewJourneyParser interface {
	mock.TestingT
	Cleanup(func())
}

// NewJourneyParser creates a new instance of JourneyParser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJourneyParser(t mockConstructorTestingTNewJourneyParser) *JourneyParser {
	mock := &JourneyParser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
