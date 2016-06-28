package mocks

import (
	"time"

	"github.com/exago/svc/repository/model"
	"github.com/stretchr/testify/mock"
)

type RepositoryData struct {
	Name, Branch string

	// Data types
	CodeStats    model.CodeStats
	Imports      model.Imports
	TestResults  model.TestResults
	LintMessages model.LintMessages
	Metadata     model.Metadata
	Score        model.Score
	StartTime    time.Time
	LastUpdate   time.Time

	mock.Mock
}

func NewRepositoryData(name, rank string) *RepositoryData {
	repoMock := RepositoryData{}
	repoMock.On("GetName").Return(name)
	repoMock.On("GetRank").Return(rank)
	repoMock.On("GetMetadataDescription").Return("metadata description")
	repoMock.On("GetMetadataImage").Return("metadata image")
	return &repoMock
}

// GetName provides a mock function with given fields:
func (_m *RepositoryData) GetName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMetadataDescription provides a mock function with given fields:
func (_m *RepositoryData) GetMetadataDescription() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMetadataImage provides a mock function with given fields:
func (_m *RepositoryData) GetMetadataImage() string {
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
func (_m *RepositoryData) GetRank() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMetadata provides a mock function with given fields:
func (_m *RepositoryData) GetMetadata() (model.Metadata, error) {
	ret := _m.Called()

	var r0 model.Metadata
	if rf, ok := ret.Get(0).(func() model.Metadata); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Metadata)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastUpdate provides a mock function with given fields:
func (_m *RepositoryData) GetLastUpdate() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecutionTime provides a mock function with given fields:
func (_m *RepositoryData) GetExecutionTime() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetScore provides a mock function with given fields:
func (_m *RepositoryData) GetScore() (model.Score, error) {
	ret := _m.Called()

	var r0 model.Score
	if rf, ok := ret.Get(0).(func() model.Score); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Score)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetImports provides a mock function with given fields:
func (_m *RepositoryData) GetImports() (model.Imports, error) {
	ret := _m.Called()

	var r0 model.Imports
	if rf, ok := ret.Get(0).(func() model.Imports); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Imports)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCodeStats provides a mock function with given fields:
func (_m *RepositoryData) GetCodeStats() (model.CodeStats, error) {
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

// GetTestResults provides a mock function with given fields:
func (_m *RepositoryData) GetTestResults() (model.TestResults, error) {
	ret := _m.Called()

	var r0 model.TestResults
	if rf, ok := ret.Get(0).(func() model.TestResults); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.TestResults)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLintMessages provides a mock function with given fields: linters
func (_m *RepositoryData) GetLintMessages(linters []string) (model.LintMessages, error) {
	ret := _m.Called(linters)

	var r0 model.LintMessages
	if rf, ok := ret.Get(0).(func([]string) model.LintMessages); ok {
		r0 = rf(linters)
	} else {
		r0 = ret.Get(0).(model.LintMessages)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(linters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetStartTime provides a mock function with given fields: t
func (_m *RepositoryData) SetStartTime(t time.Time) {
	_m.Called(t)
}

// IsCached provides a mock function with given fields:
func (_m *RepositoryData) IsCached() bool {
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
func (_m *RepositoryData) IsLoaded() bool {
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
func (_m *RepositoryData) Load() error {
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
func (_m *RepositoryData) ClearCache() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AsMap provides a mock function with given fields:
func (_m *RepositoryData) AsMap() map[string]interface{} {
	ret := _m.Called()

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}
