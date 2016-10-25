package queue

import (
	"sync"
	"testing"

	"github.com/hotolab/exago-svc/taskrunner"
)

func process(repo, branch string, tr taskrunner.TaskRunner) (interface{}, error) {
	return "foo", nil
}

func TestNewQueue(t *testing.T) {
	queue = &Queue{
		sem:  make(chan bool, 4),
		in:   make(chan string),
		out:  make(chan map[string]interface{}),
		quit: make(chan bool),
		wg:   &sync.WaitGroup{},
	}
	queue.Init()
	queue.Stop()
	if _, err := queue.PushSync("test"); err != ErrQueueIsClosing {
		t.Error("Should not be able to push items in a closing queue")
	}
	if err := queue.PushAsync("test"); err != ErrQueueIsClosing {
		t.Error("Should not be able to push items in a closing queue")
	}
}

func TestQueuePush(t *testing.T) {
	queue = &Queue{
		processorFn: process,
		sem:         make(chan bool, 4),
		in:          make(chan string, 20),
		out:         make(chan map[string]interface{}, 100),
		quit:        make(chan bool),
		wg:          &sync.WaitGroup{},
	}
	queue.Init()

	var err error

	data, err := queue.PushSync("testsync")
	if err != nil {
		t.Error("There should be no error")
	}
	if data != "foo" {
		t.Error("Unexpected output received")
	}

	if err = queue.PushAsync("testasync"); err != nil {
		t.Error("There should be no error")
	}

	queue.WaitUntilEmpty()
}
