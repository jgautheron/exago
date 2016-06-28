package showcaser

import (
	"testing"

	"github.com/exago/svc/mocks"
)

func TestGotRecentRepositories(t *testing.T) {
	data = New()

	ProcessRepository(mocks.NewRepositoryData("github.com/foo/bar", "A"))
	ProcessRepository(mocks.NewRepositoryData("github.com/bar/foo", "B"))
	ProcessRepository(mocks.NewRepositoryData("github.com/moo/foo", "D"))

	recent := GetRecentRepositories()
	if len(recent) != 3 {
		t.Errorf("Expected 3 recent repos, got %d", len(recent))
	}
}

func TestGotTopRankedRepositories(t *testing.T) {
	data = New()

	ProcessRepository(mocks.NewRepositoryData("github.com/foo/bar", "A"))
	ProcessRepository(mocks.NewRepositoryData("github.com/bar/foo", "B"))
	ProcessRepository(mocks.NewRepositoryData("github.com/moo/foo", "D"))
	ProcessRepository(mocks.NewRepositoryData("github.com/foo/boo", "A"))
	ProcessRepository(mocks.NewRepositoryData("github.com/moo/boo", "A"))
	ProcessRepository(mocks.NewRepositoryData("github.com/boo/bar", "A"))
	ProcessRepository(mocks.NewRepositoryData("github.com/bar/boo", "A"))
	ProcessRepository(mocks.NewRepositoryData("github.com/bar/bar", "A"))
	ProcessRepository(mocks.NewRepositoryData("github.com/boo/boo", "A"))

	top := GetTopRankedRepositories()
	if len(top) != ItemCount {
		t.Errorf("Expected %d top repos, got %d", ItemCount, len(top))
	}
}
