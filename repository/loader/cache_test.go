package loader

import (
	"encoding/json"
	"testing"

	. "github.com/hotolab/exago-svc"
	"github.com/hotolab/exago-svc/mocks"
	"github.com/hotolab/exago-svc/repository"
	ldb "github.com/syndtr/goleveldb/leveldb"
)

const (
	repo      = "github.com/foo/bar"
	branch    = "master"
	goversion = "1.7.4"
)

func TestDidSave(t *testing.T) {
	stub := repository.Repository{Name: repo, Branch: branch, GoVersion: goversion}
	b, _ := json.Marshal(stub)

	dbMock := mocks.Database{}
	dbMock.On("Put",
		getCacheKey(repo, branch, goversion), b,
	).Return(nil)

	rp := repository.New(repo, branch, goversion)
	l := New(
		WithDatabase(dbMock),
	)
	if err := l.Save(rp); err != nil {
		t.Errorf("Got error while saving the data: %v", err)
	}
}

func TestIsNotCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get",
		getCacheKey(repo, branch, goversion),
	).Return([]byte(""), ldb.ErrNotFound)

	l := New(
		WithDatabase(dbMock),
	)
	if cached := l.IsCached(repo, branch, goversion); cached {
		t.Errorf("The repository %s should not be cached", repo)
	}
}

func TestIsCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get",
		getCacheKey(repo, branch, goversion),
	).Return([]byte(""), nil)

	l := New(
		WithDatabase(dbMock),
	)
	if cached := l.IsCached(repo, branch, goversion); !cached {
		t.Errorf("The repository %s should be cached", repo)
	}
}

func TestCacheCleared(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Delete", []byte(
		getCacheKey(repo, branch, goversion),
	)).Return(nil)

	l := New(
		WithDatabase(dbMock),
	)
	if err := l.ClearCache(repo, branch, goversion); err != nil {
		t.Error("Got error while attempting to clear cache")
	}
}
