package showcaser

import (
	"testing"
	"time"

	. "github.com/hotolab/exago-svc/config"
	"github.com/hotolab/exago-svc/mocks"
	. "github.com/stretchr/testify/mock"
)

func init() {
	Config.ShowcaserPopularRebuildInterval = 10 * time.Millisecond
}

// Flaky test
// func TestInterrupted(t *testing.T) {
// 	dbMock := getDatabaseMock()
// 	// dbMock.On("Get", Anything).Return([]byte(repoStubData), nil)
// 	dbMock.On("Put", Anything, Anything).Return(nil)

// 	showcaser := getShowcaseMock(dbMock)
// 	go showcaser.catchInterrupt()

// 	done := make(chan bool, 1)

// 	go func() {
// 		select {
// 		case <-signals:
// 			done <- true
// 		}
// 	}()

// 	go func() {
// 		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
// 	}()

// 	select {
// 	case <-done:
// 		// Do nothing
// 	case <-time.After(200 * time.Millisecond):
// 		t.Error("Timeout waiting for SIGINT")
// 	}
// }

func TestPeriodicallyRebuilt(t *testing.T) {
	dbMock := getDatabaseMock()
	dbMock.On("Get", Anything).Return([]byte(repoStubData), nil)
	showcaser := getShowcaseMock(dbMock)

	showcaser.Process(mocks.NewRecord("github.com/moo/bar", "", "B"))

	go showcaser.periodicallyRebuildPopularList()
	time.Sleep(100 * time.Millisecond)
	showcaser.RLock()
	defer showcaser.RUnlock()
	if len(showcaser.popular) != 1 {
		t.Errorf("The popular slice should have a length of 1, got %d", len(showcaser.popular))
	}
}
