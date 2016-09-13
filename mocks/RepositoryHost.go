package mocks

import "github.com/stretchr/testify/mock"

type RepositoryHost struct {
	mock.Mock
}

// GetFileContent provides a mock function with given fields: owner, repository, path
func (_m RepositoryHost) GetFileContent(owner string, repository string, path string) (string, error) {
	ret := _m.Called(owner, repository, path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string) string); ok {
		r0 = rf(owner, repository, path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(owner, repository, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: owner, repository
func (_m RepositoryHost) Get(owner string, repository string) (map[string]interface{}, error) {
	ret := _m.Called(owner, repository)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(string, string) map[string]interface{}); ok {
		r0 = rf(owner, repository)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(owner, repository)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
