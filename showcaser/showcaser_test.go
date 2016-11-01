package showcaser

import (
	"testing"

	. "github.com/hotolab/exago-svc"
	"github.com/hotolab/exago-svc/mocks"
	"github.com/hotolab/exago-svc/repository"
	. "github.com/stretchr/testify/mock"
	ldb "github.com/syndtr/goleveldb/leveldb"
)

var (
	repoStubData = `{"metadata":{"image":"https://avatars.githubusercontent.com/u/683888?v=3","description":"A codename generator meant for naming software releases.","stars":13,"last_push":"2015-08-29T20:32:12Z"},"score":{"value":68.72365160829756,"rank":"D"}}`
	snapshotStub = `{"Recent":["github.com/foo/bar","github.com/moo/bar"],"TopRanked":["github.com/foo/bar"],"Popular":["github.com/moo/bar","github.com/foo/bar"],"Topk":"BAQA/8gO/4EEAQL/ggABDAEEAAAs/4IAAhJnaXRodWIuY29tL2Zvby9iYXIAEmdpdGh1Yi5jb20vbW9vL2JhcgIN/4UCAQL/hgAB/4QAADH/gwMBAQdFbGVtZW50Af+EAAEDAQNLZXkBDAABBUNvdW50AQQAAQVFcnJvcgEEAAAAMv+GAAIBEmdpdGh1Yi5jb20vZm9vL2JhcgECAAESZ2l0aHViLmNvbS9tb28vYmFyAQIADP+HAgEC/4gAAQQAAP4CXv+IAP4CWAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="}`
)

func getDatabaseMock() mocks.Database {
	db := mocks.Database{}
	// showcaser.New systematically looks for a snapshot
	db.On("Get", []byte(DatabaseKey)).Return(nil, ldb.ErrNotFound)
	return db
}

func getShowcaseMock(db mocks.Database) *Showcaser {
	rh := mocks.RepositoryHost{}
	rl := repository.NewLoader(
		WithDatabase(db),
		WithRepositoryHost(rh),
	)
	showcaser, _ := New(
		WithDatabase(db),
		WithRepositoryLoader(rl),
	)
	return showcaser
}

func TestNew(t *testing.T) {
	showcaser := getShowcaseMock(getDatabaseMock())
	if showcaser.itemCount != ItemCount {
		t.Error("The item count should match the default one")
	}
}

func TestSave(t *testing.T) {
	dbMock := getDatabaseMock()
	// AnythingOfType would be more appropriate but doesn't work as expected.
	// See: https://github.com/stretchr/testify/issues/68
	dbMock.On("Put", Anything, Anything).Return(nil)
	showcaser := getShowcaseMock(dbMock)
	if err := showcaser.save(); err != nil {
		t.Errorf("Got error while saving the data: %v", err)
	}
}

func TestRepositoryRankedAAdded(t *testing.T) {
	showcaser := getShowcaseMock(getDatabaseMock())
	showcaser.Process(mocks.NewRecord("github.com/foo/bar", "A"))
	if len(showcaser.topRanked) != 1 || len(showcaser.recent) != 1 || len(showcaser.tk.Keys()) != 1 {
		t.Error("There should be exactly one entry per slice")
	}
}

func TestRepositoryRankedADuplicated(t *testing.T) {
	showcaser := getShowcaseMock(getDatabaseMock())
	showcaser.Process(mocks.NewRecord("github.com/foo/bar", "A"))
	showcaser.Process(mocks.NewRecord("github.com/foo/bar", "A"))
	if len(showcaser.topRanked) != 1 || len(showcaser.recent) != 1 || len(showcaser.tk.Keys()) != 1 {
		t.Error("There should be exactly one entry per slice")
	}
}

func TestRepositoryRankedBAdded(t *testing.T) {
	showcaser := getShowcaseMock(getDatabaseMock())
	showcaser.Process(mocks.NewRecord("github.com/moo/bar", "B"))
	if len(showcaser.topRanked) != 0 || len(showcaser.recent) != 1 || len(showcaser.tk.Keys()) != 1 {
		t.Error("There should be exactly one entry per slice")
	}
}

func TestDataSerialized(t *testing.T) {
	showcaser := getShowcaseMock(getDatabaseMock())
	showcaser.Process(mocks.NewRecord("github.com/moo/bar", "B"))
	_, err := showcaser.serialize()
	if err != nil {
		t.Errorf("The serialization went wrong: %v", err)
	}
}

func TestPopularDataUpdated(t *testing.T) {
	dbMock := getDatabaseMock()
	dbMock.On("Get", Anything).Return([]byte(repoStubData), nil)
	showcaser := getShowcaseMock(dbMock)

	showcaser.Process(mocks.NewRecord("github.com/moo/bar", "B"))
	err := showcaser.updatePopular()
	if err != nil {
		t.Errorf("Everything should go fine: %v", err)
	}
	if len(showcaser.popular) != 1 {
		t.Errorf("The popular slice should have a length of 1, got %d", len(showcaser.popular))
	}
}

func TestDataLoadedFromDB(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get", []byte(DatabaseKey)).Return([]byte(snapshotStub), nil)
	dbMock.On("Get", Anything).Return([]byte(repoStubData), nil)
	showcaser := getShowcaseMock(dbMock)
	_, exists, err := showcaser.loadFromDB()
	if !exists {
		t.Error("The snapshot should exist")
	}
	if err != nil {
		t.Errorf("Got error while loading the snapshot from DB: %v", err)
	}
}
