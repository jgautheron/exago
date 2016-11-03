package loader_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/hotolab/exago-svc"
	"github.com/hotolab/exago-svc/mocks"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/repository/loader"
	"github.com/hotolab/exago-svc/repository/model"
	ldb "github.com/syndtr/goleveldb/leveldb"
)

const (
	repo   = "test"
	branch = ""
)

func TestDidSave(t *testing.T) {
	stub := model.Data{}
	b, _ := json.Marshal(stub)

	dbMock := mocks.Database{}
	dbMock.On("Put",
		[]byte(fmt.Sprintf("%s-%s", repo, branch)), b,
	).Return(nil)

	rp := repository.New(repo, branch)
	l := loader.New(
		WithDatabase(dbMock),
	)
	if err := l.Save(rp); err != nil {
		t.Errorf("Got error while saving the data: %v", err)
	}
}

func TestIsNotCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get",
		[]byte(fmt.Sprintf("%s-%s", repo, branch)),
	).Return([]byte(""), ldb.ErrNotFound)

	l := loader.New(
		WithDatabase(dbMock),
	)
	if cached := l.IsCached(repo, branch); cached {
		t.Errorf("The repository %s should not be cached", repo)
	}
}

func TestIsCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get",
		[]byte(fmt.Sprintf("%s-%s", repo, branch)),
	).Return([]byte(""), nil)

	l := loader.New(
		WithDatabase(dbMock),
	)
	if cached := l.IsCached(repo, branch); !cached {
		t.Errorf("The repository %s should be cached", repo)
	}
}

func TestCacheCleared(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Delete", []byte(
		fmt.Sprintf("%s-%s", repo, branch),
	)).Return(nil)

	l := loader.New(
		WithDatabase(dbMock),
	)
	if err := l.ClearCache(repo, branch); err != nil {
		t.Error("Got error while attempting to clear cache")
	}
}
