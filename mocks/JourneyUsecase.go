// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	io "io"

	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// JourneyUsecase is an autogenerated mock type for the JourneyUsecase type
type JourneyUsecase struct {
	mock.Mock
}

// ImportFromCSVFile provides a mock function with given fields: c, reader
func (_m *JourneyUsecase) ImportFromCSVFile(c *gin.Context, reader io.Reader) (int64, []string) {
	ret := _m.Called(c, reader)

	var r0 int64
	var r1 []string
	if rf, ok := ret.Get(0).(func(*gin.Context, io.Reader) (int64, []string)); ok {
		return rf(c, reader)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context, io.Reader) int64); ok {
		r0 = rf(c, reader)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(*gin.Context, io.Reader) []string); ok {
		r1 = rf(c, reader)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewJourneyUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewJourneyUsecase creates a new instance of JourneyUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJourneyUsecase(t mockConstructorTestingTNewJourneyUsecase) *JourneyUsecase {
	mock := &JourneyUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
