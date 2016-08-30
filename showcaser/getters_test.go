package showcaser

import (
	"testing"

	"github.com/hotolab/exago-svc/mocks"
	"github.com/hotolab/exago-svc/repository"
	. "github.com/stretchr/testify/mock"
)

func getRecordMock(repo, rank string) repository.Record {
	mock := mocks.NewRecord(repo, rank)
	return mock
}

func TestGotRecentRepositories(t *testing.T) {
	showcase = getShowcaseMock(mocks.Database{})

	showcase.Process(getRecordMock("github.com/foo/bar", "A"))
	showcase.Process(getRecordMock("github.com/bar/foo", "B"))
	showcase.Process(getRecordMock("github.com/moo/foo", "D"))

	recent := showcase.GetRecentRepositories()
	if len(recent) != 3 {
		t.Errorf("Expected 3 recent repos, got %d", len(recent))
	}
}

func TestGotTopRankedRepositories(t *testing.T) {
	showcase = getShowcaseMock(mocks.Database{})

	showcase.Process(getRecordMock("github.com/foo/bar", "A"))
	showcase.Process(getRecordMock("github.com/bar/foo", "B"))
	showcase.Process(getRecordMock("github.com/moo/foo", "D"))
	showcase.Process(getRecordMock("github.com/foo/boo", "A"))
	showcase.Process(getRecordMock("github.com/moo/boo", "A"))
	showcase.Process(getRecordMock("github.com/boo/bar", "A"))
	showcase.Process(getRecordMock("github.com/bar/boo", "A"))
	showcase.Process(getRecordMock("github.com/bar/bar", "A"))
	showcase.Process(getRecordMock("github.com/boo/boo", "A"))

	top := showcase.GetTopRankedRepositories()
	if len(top) != ItemCount {
		t.Errorf("Expected %d top repos, got %d", ItemCount, len(top))
	}
}

func TestGotPopularRepositories(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get", Anything).Return([]byte(repoStubData), nil)
	showcase = getShowcaseMock(dbMock)

	showcase.Process(getRecordMock("github.com/foo/bar", "A"))
	showcase.Process(getRecordMock("github.com/foo/bar", "A"))
	showcase.Process(getRecordMock("github.com/bar/foo", "B"))
	showcase.Process(getRecordMock("github.com/moo/foo", "D"))
	showcase.Process(getRecordMock("github.com/foo/boo", "A"))
	showcase.Process(getRecordMock("github.com/moo/boo", "A"))
	showcase.Process(getRecordMock("github.com/boo/bar", "A"))
	showcase.Process(getRecordMock("github.com/bar/boo", "A"))
	showcase.Process(getRecordMock("github.com/bar/boo", "A"))
	showcase.Process(getRecordMock("github.com/bar/boo", "A"))
	showcase.Process(getRecordMock("github.com/bar/bar", "A"))
	showcase.Process(getRecordMock("github.com/boo/boo", "A"))

	if err := showcase.updatePopular(); err != nil {
		t.Errorf("Got error while updating the popular list: %v", err)
	}

	popular := showcase.GetPopularRepositories()
	if len(popular) != ItemCount {
		t.Errorf("Expected %d popular repos, got %d", ItemCount, len(popular))
	}
}
