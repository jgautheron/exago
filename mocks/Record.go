package mocks

import "github.com/stretchr/testify/mock"

import "time"
import "github.com/exago/svc/repository/model"

type Record struct {
	mock.Mock
}

func NewRecord(name, rank string) *Record {
	repoMock := Record{}
	repoMock.On("GetName").Return(name)
	repoMock.On("GetRank").Return(rank)
	repoMock.On("GetMetadataDescription").Return("metadata description")
	repoMock.On("GetMetadataImage").Return("metadata image")
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

// SetMetadata provides a mock function with given fields:
func (_m *Record) SetMetadata() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
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

// SetLastUpdate provides a mock function with given fields:
func (_m *Record) SetLastUpdate() {
	_m.Called()
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

// SetExecutionTime provides a mock function with given fields:
func (_m *Record) SetExecutionTime() {
	_m.Called()
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

// SetScore provides a mock function with given fields:
func (_m *Record) SetScore() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetImports provides a mock function with given fields:
func (_m *Record) GetImports() model.Imports {
	ret := _m.Called()

	var r0 model.Imports
	if rf, ok := ret.Get(0).(func() model.Imports); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Imports)
	}

	return r0
}

// SetImports provides a mock function with given fields: _a0
func (_m *Record) SetImports(_a0 model.Imports) {
	_m.Called(_a0)
}

// GetCodeStats provides a mock function with given fields:
func (_m *Record) GetCodeStats() model.CodeStats {
	ret := _m.Called()

	var r0 model.CodeStats
	if rf, ok := ret.Get(0).(func() model.CodeStats); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.CodeStats)
	}

	return r0
}

// SetCodeStats provides a mock function with given fields: _a0
func (_m *Record) SetCodeStats(_a0 model.CodeStats) {
	_m.Called(_a0)
}

// GetTestResults provides a mock function with given fields:
func (_m *Record) GetTestResults() model.TestResults {
	ret := _m.Called()

	var r0 model.TestResults
	if rf, ok := ret.Get(0).(func() model.TestResults); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.TestResults)
	}

	return r0
}

// SetTestResults provides a mock function with given fields: tr
func (_m *Record) SetTestResults(tr model.TestResults) {
	_m.Called(tr)
}

// GetLintMessages provides a mock function with given fields: linters
func (_m *Record) GetLintMessages(linters []string) model.LintMessages {
	ret := _m.Called(linters)

	var r0 model.LintMessages
	if rf, ok := ret.Get(0).(func([]string) model.LintMessages); ok {
		r0 = rf(linters)
	} else {
		r0 = ret.Get(0).(model.LintMessages)
	}

	return r0
}

// SetLintMessages provides a mock function with given fields: _a0
func (_m *Record) SetLintMessages(_a0 model.LintMessages) {
	_m.Called(_a0)
}

// SetStartTime provides a mock function with given fields: t
func (_m *Record) SetStartTime(t time.Time) {
	_m.Called(t)
}

// SetError provides a mock function with given fields: tp, err
func (_m *Record) SetError(tp string, err error) {
	_m.Called(tp, err)
}

// IsCached provides a mock function with given fields:
func (_m *Record) IsCached() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsLoaded provides a mock function with given fields:
func (_m *Record) IsLoaded() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Load provides a mock function with given fields:
func (_m *Record) Load() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearCache provides a mock function with given fields:
func (_m *Record) ClearCache() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields:
func (_m *Record) Save() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
