package showcaser

import (
	"testing"

	"github.com/exago/svc/mocks"
)

func TestGotRecentRepositories(t *testing.T) {
	data = New()

	ProcessRepository(mocks.NewRepositoryMock("github.com/foo/bar", "A"))
	ProcessRepository(mocks.NewRepositoryMock("github.com/bar/foo", "B"))
	ProcessRepository(mocks.NewRepositoryMock("github.com/moo/foo", "D"))

	recent := GetRecentRepositories()
	if len(recent) != 3 {
		t.Errorf("Expected 3 recent repos, got %d", len(recent))
	}
}

func TestGotTopRankedRepositories(t *testing.T) {
	data = New()

	ProcessRepository(mocks.NewRepositoryMock("github.com/foo/bar", "A"))
	ProcessRepository(mocks.NewRepositoryMock("github.com/bar/foo", "B"))
	ProcessRepository(mocks.NewRepositoryMock("github.com/moo/foo", "D"))
	ProcessRepository(mocks.NewRepositoryMock("github.com/bar/boo", "A"))

	top := GetTopRankedRepositories()
	if len(top) != 2 {
		t.Errorf("Expected 2 top repos, got %d", len(top))
	}
}
