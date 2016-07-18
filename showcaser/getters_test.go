package showcaser

import (
	"testing"

	"github.com/exago/svc/mocks"
	"github.com/exago/svc/repository"
)

func getRecordMock(repo, rank string) repository.Record {
	mock := mocks.NewRecord(repo, rank)
	return mock
}

func TestGotRecentRepositories(t *testing.T) {
	data = getShowcaseMock()

	ProcessRepository(getRecordMock("github.com/foo/bar", "A"))
	ProcessRepository(getRecordMock("github.com/bar/foo", "B"))
	ProcessRepository(getRecordMock("github.com/moo/foo", "D"))

	recent := GetRecentRepositories()
	if len(recent) != 3 {
		t.Errorf("Expected 3 recent repos, got %d", len(recent))
	}
}

func TestGotTopRankedRepositories(t *testing.T) {
	data = getShowcaseMock()

	ProcessRepository(getRecordMock("github.com/foo/bar", "A"))
	ProcessRepository(getRecordMock("github.com/bar/foo", "B"))
	ProcessRepository(getRecordMock("github.com/moo/foo", "D"))
	ProcessRepository(getRecordMock("github.com/foo/boo", "A"))
	ProcessRepository(getRecordMock("github.com/moo/boo", "A"))
	ProcessRepository(getRecordMock("github.com/boo/bar", "A"))
	ProcessRepository(getRecordMock("github.com/bar/boo", "A"))
	ProcessRepository(getRecordMock("github.com/bar/bar", "A"))
	ProcessRepository(getRecordMock("github.com/boo/boo", "A"))

	top := GetTopRankedRepositories()
	if len(top) != ItemCount {
		t.Errorf("Expected %d top repos, got %d", ItemCount, len(top))
	}
}
