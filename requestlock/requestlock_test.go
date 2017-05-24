package requestlock_test

import (
	"testing"

	"github.com/jgautheron/exago/requestlock"
)

func TestLockAdded(t *testing.T) {
	requestlock.Add("127.0.0.1", "github.com/foo/bar")
	if !requestlock.Contains("127.0.0.1", "github.com/foo/bar") {
		t.Error("Could not find the lock")
	}
}

func TestLockRemoved(t *testing.T) {
	requestlock.Remove("127.0.0.1", "github.com/foo/bar")
	if requestlock.Contains("127.0.0.1", "github.com/foo/bar") {
		t.Error("The lock has not been removed")
	}
}
