package mocks

import "github.com/hotolab/exago-svc/repository/model"
import "github.com/stretchr/testify/mock"

import "time"

type Record struct {
	mock.Mock
}

func NewRecord(name, branch, rank string) *Record {
	repoMock := Record{}
	repoMock.On("GetName").Return(name)
	repoMock.On("GetBranch").Return(branch)
	repoMock.On("GetRank").Return(rank)
	repoMock.On("GetMetadata").Return(model.Metadata{
		Image:       "http://foo.com/img.png",
		Description: "This is a repository description",
		Stars:       912,
		LastPush:    time.Now(),
	})
	return &repoMock
}

// GetName provides a mock function with given fields:
func (_m *Record) GetName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetBranch provides a mock function with given fields:
func (_m *Record) GetBranch() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetName provides a mock function with given fields: _a0
func (_m *Record) SetName(_a0 string) {
	_m.Called(_a0)
}

// GetRank provides a mock function with given fields:
func (_m *Record) GetRank() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetData provides a mock function with given fields:
func (_m *Record) GetData() model.Data {
	ret := _m.Called()

	var r0 model.Data
	if rf, ok := ret.Get(0).(func() model.Data); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Data)
	}

	return r0
}

// SetData provides a mock function with given fields: d
func (_m *Record) SetData(d model.Data) {
	_m.Called(d)
}

// GetMetadata provides a mock function with given fields:
func (_m *Record) GetMetadata() model.Metadata {
	ret := _m.Called()

	var r0 model.Metadata
	if rf, ok := ret.Get(0).(func() model.Metadata); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Metadata)
	}

	return r0
}

// SetMetadata provides a mock function with given fields: _a0
func (_m *Record) SetMetadata(_a0 model.Metadata) {
	_m.Called(_a0)
}

// GetLastUpdate provides a mock function with given fields:
func (_m *Record) GetLastUpdate() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// SetLastUpdate provides a mock function with given fields: t
func (_m *Record) SetLastUpdate(t time.Time) {
	_m.Called(t)
}

// GetExecutionTime provides a mock function with given fields:
func (_m *Record) GetExecutionTime() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetExecutionTime provides a mock function with given fields: duration
func (_m *Record) SetExecutionTime(duration time.Duration) {
	_m.Called(duration)
}

// GetScore provides a mock function with given fields:
func (_m *Record) GetScore() model.Score {
	ret := _m.Called()

	var r0 model.Score
	if rf, ok := ret.Get(0).(func() model.Score); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Score)
	}

	return r0
}

// GetProjectRunner provides a mock function with given fields:
func (_m *Record) GetProjectRunner() model.ProjectRunner {
	ret := _m.Called()

	var r0 model.ProjectRunner
	if rf, ok := ret.Get(0).(func() model.ProjectRunner); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.ProjectRunner)
	}

	return r0
}

// SetProjectRunner provides a mock function with given fields: tr
func (_m *Record) SetProjectRunner(tr model.ProjectRunner) {
	_m.Called(tr)
}

// GetLintMessages provides a mock function with given fields:
func (_m *Record) GetLintMessages() model.LintMessages {
	ret := _m.Called()

	var r0 model.LintMessages
	if rf, ok := ret.Get(0).(func() model.LintMessages); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.LintMessages)
	}

	return r0
}

// SetLintMessages provides a mock function with given fields: _a0
func (_m *Record) SetLintMessages(_a0 model.LintMessages) {
	_m.Called(_a0)
}

// SetError provides a mock function with given fields: tp, err
func (_m *Record) SetError(tp string, err error) {
	_m.Called(tp, err)
}

// HasError provides a mock function with given fields:
func (_m *Record) HasError() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ApplyScore provides a mock function with given fields:
func (_m *Record) ApplyScore() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
