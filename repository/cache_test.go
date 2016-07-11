package repository

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/exago/svc/mocks"
	ldb "github.com/syndtr/goleveldb/leveldb"
)

func TestDidSave(t *testing.T) {
	stub, err := loadStubRepo()
	if err != nil {
		t.Error(err)
	}

	dbMock := mocks.Database{}
	dbMock.On("Get",
		[]byte(fmt.Sprintf("%s-%s", repo, "")),
	).Return([]byte(data), nil)

	b, _ := json.Marshal(stub.Data)
	dbMock.On("Put",
		[]byte(fmt.Sprintf("%s-%s", repo, "")), b,
	).Return(nil)

	rp := &Repository{
		name: repo,
		db:   dbMock,
	}
	rp.Load()
	if err := rp.Save(); err != nil {
		t.Errorf("Got error while saving the data: %v", err)
	}
}

func TestIsNotCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get",
		[]byte(fmt.Sprintf("%s-%s", repo, "")),
	).Return([]byte(""), ldb.ErrNotFound)

	rp := &Repository{
		name: repo,
		db:   dbMock,
	}
	cached := rp.IsCached()
	if cached {
		t.Errorf("The repository %s should not be cached", rp.name)
	}
}

func TestIsCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get",
		[]byte(fmt.Sprintf("%s-%s", repo, "")),
	).Return([]byte(""), nil)

	rp := &Repository{
		name: repo,
		db:   dbMock,
	}
	cached := rp.IsCached()
	if !cached {
		t.Errorf("The repository %s should be cached", rp.name)
	}
}

func TestCacheCleared(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Delete", []byte(
		fmt.Sprintf("%s-%s", repo, ""),
	)).Return(nil)

	rp := &Repository{
		db:   dbMock,
		name: repo,
	}
	if err := rp.ClearCache(); err != nil {
		t.Error("Got error while attempting to clear cache")
	}
}
