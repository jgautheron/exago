package mocks

import (
	"github.com/hotolab/exago-svc/repository/model"
	"github.com/stretchr/testify/mock"
)

type TaskRunner struct {
	mock.Mock
}

// FetchCodeStats provides a mock function with given fields:
func (_m TaskRunner) FetchCodeStats() (model.CodeStats, error) {
	ret := _m.Called()

	var r0 model.CodeStats
	if rf, ok := ret.Get(0).(func() model.CodeStats); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.CodeStats)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchLintMessages provides a mock function with given fields: linters
func (_m TaskRunner) FetchLintMessages() (model.LintMessages, error) {
	ret := _m.Called()

	var r0 model.LintMessages
	if rf, ok := ret.Get(0).(func() model.LintMessages); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.LintMessages)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchProjectRunner provides a mock function with given fields:
func (_m TaskRunner) FetchProjectRunner() (model.ProjectRunner, error) {
	ret := _m.Called()

	var r0 model.ProjectRunner
	if rf, ok := ret.Get(0).(func() model.ProjectRunner); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.ProjectRunner)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
