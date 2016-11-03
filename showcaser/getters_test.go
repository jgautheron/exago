package showcaser

import (
	"testing"

	"github.com/hotolab/exago-svc/mocks"
	"github.com/hotolab/exago-svc/repository/model"
	. "github.com/stretchr/testify/mock"
)

func getRecordMock(repo, rank string) model.Record {
	mock := mocks.NewRecord(repo, "", rank)
	return mock
}

func TestGotRecentRepositories(t *testing.T) {
	showcaser := getShowcaseMock(getDatabaseMock())

	showcaser.Process(getRecordMock("github.com/foo/bar", "A"))
	showcaser.Process(getRecordMock("github.com/bar/foo", "B"))
	showcaser.Process(getRecordMock("github.com/moo/foo", "D"))

	recent := showcaser.GetRecentRepositories()
	if len(recent) != 3 {
		t.Errorf("Expected 3 recent repos, got %d", len(recent))
	}
}

func TestGotTopRankedRepositories(t *testing.T) {
	showcaser := getShowcaseMock(getDatabaseMock())

	showcaser.Process(getRecordMock("github.com/foo/bar", "A"))
	showcaser.Process(getRecordMock("github.com/bar/foo", "B"))
	showcaser.Process(getRecordMock("github.com/moo/foo", "D"))
	showcaser.Process(getRecordMock("github.com/foo/boo", "A"))
	showcaser.Process(getRecordMock("github.com/moo/boo", "A"))
	showcaser.Process(getRecordMock("github.com/boo/bar", "A"))
	showcaser.Process(getRecordMock("github.com/bar/boo", "A"))
	showcaser.Process(getRecordMock("github.com/bar/bar", "A"))
	showcaser.Process(getRecordMock("github.com/boo/boo", "A"))

	top := showcaser.GetTopRankedRepositories()
	if len(top) != ItemCount {
		t.Errorf("Expected %d top repos, got %d", ItemCount, len(top))
	}
}

func TestGotPopularRepositories(t *testing.T) {
	dbMock := getDatabaseMock()
	dbMock.On("Get", Anything).Return([]byte(repoStubData), nil)
	showcaser := getShowcaseMock(dbMock)

	showcaser.Process(getRecordMock("github.com/foo/bar", "A"))
	showcaser.Process(getRecordMock("github.com/foo/bar", "A"))
	showcaser.Process(getRecordMock("github.com/bar/foo", "B"))
	showcaser.Process(getRecordMock("github.com/moo/foo", "D"))
	showcaser.Process(getRecordMock("github.com/foo/boo", "A"))
	showcaser.Process(getRecordMock("github.com/moo/boo", "A"))
	showcaser.Process(getRecordMock("github.com/boo/bar", "A"))
	showcaser.Process(getRecordMock("github.com/bar/boo", "A"))
	showcaser.Process(getRecordMock("github.com/bar/boo", "A"))
	showcaser.Process(getRecordMock("github.com/bar/boo", "A"))
	showcaser.Process(getRecordMock("github.com/bar/bar", "A"))
	showcaser.Process(getRecordMock("github.com/boo/boo", "A"))

	if err := showcaser.updatePopular(); err != nil {
		t.Errorf("Got error while updating the popular list: %v", err)
	}

	popular := showcaser.GetPopularRepositories()
	if len(popular) != ItemCount {
		t.Errorf("Expected %d popular repos, got %d", ItemCount, len(popular))
	}
}
