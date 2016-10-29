package showcaser

import (
	"syscall"
	"testing"
	"time"

	. "github.com/hotolab/exago-svc/config"
	"github.com/hotolab/exago-svc/mocks"
	. "github.com/stretchr/testify/mock"
)

func init() {
	Config.ShowcaserPopularRebuildInterval = 10 * time.Millisecond
}

func TestInterrupted(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Put", Anything, Anything).Return(nil)
	showcase = getShowcaseMock(dbMock)

	done := make(chan bool, 1)
	go showcase.catchInterrupt()

	go func() {
		select {
		case <-signals:
			done <- true
		}
	}()

	go func() {
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	select {
	case <-done:
		// Do nothing
	case <-time.After(200 * time.Millisecond):
		t.Error("Timeout waiting for SIGINT")
	}
}

func TestPeriodicallyRebuilt(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("Get", Anything).Return([]byte(repoStubData), nil)
	showcase = getShowcaseMock(dbMock)

	showcase.Process(mocks.NewRecord("github.com/moo/bar", "B"))

	go showcase.periodicallyRebuildPopularList()
	time.Sleep(25 * time.Millisecond)
	showcase.RLock()
	defer showcase.RUnlock()
	if len(showcase.popular) != 1 {
		t.Errorf("The popular slice should have a length of 1, got %d", len(showcase.popular))
	}
}
