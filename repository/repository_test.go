package repository

import (
	"fmt"
	"testing"

	"github.com/exago/svc/mocks"
	"github.com/exago/svc/repository/model"
)

var repo = "github.com/foo/bar"

func TestIsNotCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("FindAllForRepository", []byte(
		fmt.Sprintf("%s-%s", repo, "")),
	).Return(map[string][]byte{}, nil)

	rp := &Repository{
		Name: repo,
		db:   dbMock,
	}
	cached := rp.IsCached()
	if cached {
		t.Errorf("The repository %s should not be cached", rp.Name)
	}
}

func TestIsCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("FindAllForRepository", []byte(
		fmt.Sprintf("%s-%s", repo, "")),
	).Return(map[string][]byte{
		"codestats":      []byte(""),
		"imports":        []byte(""),
		"testresults":    []byte(""),
		"lintmessages":   []byte(""),
		"metadata":       []byte(""),
		"score":          []byte(""),
		"execution_time": []byte(""),
		"last_update":    []byte(""),
	}, nil)

	rp := &Repository{
		Name: repo,
		db:   dbMock,
	}
	cached := rp.IsCached()
	if !cached {
		t.Errorf("The repository %s should be cached", rp.Name)
	}
}

func TestIsNotLoaded(t *testing.T) {
	rp := &Repository{
		Name: repo,
	}
	loaded := rp.IsLoaded()
	if loaded {
		t.Errorf("The repository %s should not be loaded", rp.Name)
	}
}

func TestIsLoaded(t *testing.T) {
	tr := model.TestResults{}
	tr.RawOutput.Gotest = "foo"

	rp := &Repository{
		Name:         repo,
		CodeStats:    model.CodeStats{"LOC": 10},
		Imports:      []string{"foo"},
		TestResults:  tr,
		LintMessages: model.LintMessages{},
	}
	loaded := rp.IsLoaded()
	if !loaded {
		t.Errorf("The repository %s should be loaded", rp.Name)
	}
}
