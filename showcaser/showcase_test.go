package showcaser

import (
	"log"
	"testing"

	"github.com/dgryski/go-topk"
	"github.com/exago/svc/mocks"
	. "github.com/stretchr/testify/mock"
)

func getShowcaseMock() Showcase {
	dbMock := mocks.Database{}

	// AnythingOfType would be more appropriate but doesn't work as expected.
	// See: https://github.com/stretchr/testify/issues/68
	dbMock.On("Put", Anything, Anything).Return(nil)

	return Showcase{
		itemCount: ItemCount,
		tk:        topk.New(TopkCount),
		db:        dbMock,
	}
}

func TestNew(t *testing.T) {
	data = getShowcaseMock()
	if data.itemCount != ItemCount {
		t.Error("The item count should match the default one")
	}
}

func TestSave(t *testing.T) {
	data = getShowcaseMock()
	if err := data.save(); err != nil {
		t.Errorf("Got error while saving the data: %v", err)
	}
}

func TestRepositoryRankedAAdded(t *testing.T) {
	data = getShowcaseMock()
	ProcessRepository(mocks.NewRecord("github.com/foo/bar", "A"))
	if len(data.topRanked) != 1 || len(data.recent) != 1 || len(data.tk.Keys()) != 1 {
		t.Error("There should be exactly one entry per slice")
	}
}

func TestRepositoryRankedADuplicated(t *testing.T) {
	data = getShowcaseMock()
	ProcessRepository(mocks.NewRecord("github.com/foo/bar", "A"))
	ProcessRepository(mocks.NewRecord("github.com/foo/bar", "A"))
	if len(data.topRanked) != 1 || len(data.recent) != 1 || len(data.tk.Keys()) != 1 {
		t.Error("There should be exactly one entry per slice")
	}
}

func TestRepositoryRankedBAdded(t *testing.T) {
	data = getShowcaseMock()
	ProcessRepository(mocks.NewRecord("github.com/moo/bar", "B"))
	if len(data.topRanked) != 0 || len(data.recent) != 1 || len(data.tk.Keys()) != 1 {
		log.Println(len(data.topRanked), len(data.recent), len(data.tk.Keys()))
		t.Error("There should be exactly one entry per slice")
	}
}

func TestDataSerialized(t *testing.T) {
	data = getShowcaseMock()
	ProcessRepository(mocks.NewRecord("github.com/moo/bar", "B"))
	_, err := data.serialize()
	if err != nil {
		t.Errorf("The serialization went wrong: %v", err)
	}
}

// func TestPopularDataUpdated(t *testing.T) {
// 	err := data.updatePopular()
// 	if err != nil {
// 		t.Errorf("Everything should go fine: %v", err)
// 	}
// 	if len(data.popular) != 1 {
// 		t.Errorf("The popular slice should have a length of 1, got %d", len(data.popular))
// 	}
// }
