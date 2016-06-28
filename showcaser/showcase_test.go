package showcaser

import (
	"log"
	"testing"

	"github.com/exago/svc/mocks"
)

func TestNew(t *testing.T) {
	data = New()
	if data.itemCount != ItemCount {
		t.Error("The item count should match the default one")
	}
}

func TestRepositoryRankedAAdded(t *testing.T) {
	data = New()
	ProcessRepository(mocks.NewRepositoryData("github.com/foo/bar", "A"))
	if len(data.topRanked) != 1 || len(data.recent) != 1 || len(data.tk.Keys()) != 1 {
		t.Error("There should be exactly one entry per slice")
	}
}

func TestRepositoryRankedADuplicated(t *testing.T) {
	data = New()
	ProcessRepository(mocks.NewRepositoryData("github.com/foo/bar", "A"))
	ProcessRepository(mocks.NewRepositoryData("github.com/foo/bar", "A"))
	if len(data.topRanked) != 1 || len(data.recent) != 1 || len(data.tk.Keys()) != 1 {
		t.Error("There should be exactly one entry per slice")
	}
}

func TestRepositoryRankedBAdded(t *testing.T) {
	data = New()
	ProcessRepository(mocks.NewRepositoryData("github.com/moo/bar", "B"))
	if len(data.topRanked) != 0 || len(data.recent) != 1 || len(data.tk.Keys()) != 1 {
		log.Println(len(data.topRanked), len(data.recent), len(data.tk.Keys()))
		t.Error("There should be exactly one entry per slice")
	}
}

func TestDataSerialized(t *testing.T) {
	data = New()
	ProcessRepository(mocks.NewRepositoryData("github.com/moo/bar", "B"))
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
