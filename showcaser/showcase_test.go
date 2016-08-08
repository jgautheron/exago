package showcaser

import (
	"log"
	"testing"

	"github.com/dgryski/go-topk"
	"github.com/hotolab/exago-svc/mocks"
	. "github.com/stretchr/testify/mock"
)

var (
	repoStubData = `{"metadata":{"image":"https://avatars.githubusercontent.com/u/683888?v=3","description":"A codename generator meant for naming software releases.","stars":13,"last_push":"2015-08-29T20:32:12Z"},"score":{"value":68.72365160829756,"rank":"D"}}`
	snapshotStub = `{"Recent":["github.com/foo/bar","github.com/moo/bar"],"TopRanked":["github.com/foo/bar"],"Popular":["github.com/moo/bar","github.com/foo/bar"],"Topk":"BAQA/8gO/4EEAQL/ggABDAEEAAAs/4IAAhJnaXRodWIuY29tL2Zvby9iYXIAEmdpdGh1Yi5jb20vbW9vL2JhcgIN/4UCAQL/hgAB/4QAADH/gwMBAQdFbGVtZW50Af+EAAEDAQNLZXkBDAABBUNvdW50AQQAAQVFcnJvcgEEAAAAMv+GAAIBEmdpdGh1Yi5jb20vZm9vL2JhcgECAAESZ2l0aHViLmNvbS9tb28vYmFyAQIADP+HAgEC/4gAAQQAAP4CXv+IAP4CWAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="}`
)

func getShowcaseMock(m mocks.Database) *Showcase {
	return &Showcase{
		itemCount: ItemCount,
		tk:        topk.New(TopkCount),
		db:        m,
	}
}

func TestNew(t *testing.T) {
	showcase = getShowcaseMock(mocks.Database{})
	if showcase.itemCount != ItemCount {
		t.Error("The item count should match the default one")
	}
}

func TestSave(t *testing.T) {
	dbMock := mocks.Database{}
	// AnythingOfType would be more appropriate but doesn't work as expected.
	// See: https://github.com/stretchr/testify/issues/68
	dbMock.On("Put", Anything, Anything).Return(nil)
	showcase = getShowcaseMock(dbMock)
	if err := showcase.save(); err != nil {
		t.Errorf("Got error while saving the data: %v", err)
	}
}

func TestRepositoryRankedAAdded(t *testing.T) {
	showcase = getShowcaseMock(mocks.Database{})
	showcase.Process(mocks.NewRecord("github.com/foo/bar", "A"))
	if len(showcase.topRanked) != 1 || len(showcase.recent) != 1 || len(showcase.tk.Keys()) != 1 {
		t.Error("There should be exactly one entry per slice")
	}
}

func TestRepositoryRankedADuplicated(t *testing.T) {
	showcase = getShowcaseMock(mocks.Database{})
	showcase.Process(mocks.NewRecord("github.com/foo/bar", "A"))
	showcase.Process(mocks.NewRecord("github.com/foo/bar", "A"))
	if len(showcase.topRanked) != 1 || len(showcase.recent) != 1 || len(showcase.tk.Keys()) != 1 {
		t.Error("There should be exactly one entry per slice")
	}
}

func TestRepositoryRankedBAdded(t *testing.T) {
	showcase = getShowcaseMock(mocks.Database{})
	showcase.Process(mocks.NewRecord("github.com/moo/bar", "B"))
	if len(showcase.topRanked) != 0 || len(showcase.recent) != 1 || len(showcase.tk.Keys()) != 1 {
		log.Println(len(showcase.topRanked), len(showcase.recent), len(showcase.tk.Keys()))
		t.Error("There should be exactly one entry per slice")
	}
}

func TestDataSerialized(t *testing.T) {
	showcase = getShowcaseMock(mocks.Database{})
	showcase.Process(mocks.NewRecord("github.com/moo/bar", "B"))
	_, err := showcase.serialize()
	if err != nil {
		t.Errorf("The serialization went wrong: %v", err)
	}
}

func TestPopularDataUpdated(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get", Anything).Return([]byte(repoStubData), nil)
	showcase = getShowcaseMock(dbMock)

	showcase.Process(mocks.NewRecord("github.com/moo/bar", "B"))
	err := showcase.updatePopular()
	if err != nil {
		t.Errorf("Everything should go fine: %v", err)
	}
	if len(showcase.popular) != 1 {
		t.Errorf("The popular slice should have a length of 1, got %d", len(showcase.popular))
	}
}

func TestDataLoadedFromDB(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get", []byte(DatabaseKey)).Return([]byte(snapshotStub), nil)
	dbMock.On("Get", Anything).Return([]byte(repoStubData), nil)
	showcase = getShowcaseMock(dbMock)
	_, exists, err := showcase.loadFromDB()
	if !exists {
		t.Error("The snapshot should exist")
	}
	if err != nil {
		t.Errorf("Got error while loading the snapshot from DB: %v", err)
	}
}
