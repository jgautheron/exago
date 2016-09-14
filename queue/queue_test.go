package queue

import (
	"sync"
	"testing"

	"github.com/hotolab/exago-svc/taskrunner"
)

func process(repo string, tr taskrunner.TaskRunner) (interface{}, error) {
	return "foo", nil
}

func TestNewQueue(t *testing.T) {
	queue = &PriorityQueue{
		concurrency: 4,
		in:          make(chan *Item, 1000),
		out:         make(chan map[uint32]interface{}, 100),
		quit:        make(chan bool),
		ready:       make(chan bool),
		wg:          &sync.WaitGroup{},
	}
	queue.Init()
	queue.WaitUntilReady()
	if len(queue.AvailableWorkers()) != 4 {
		t.Error("The queue should have 4 workers")
	}
	queue.Stop()
	if len(queue.AvailableWorkers()) != 0 {
		t.Error("All workers should be closed")
	}
	if _, err := queue.PushSync("test", 10); err != ErrQueueIsClosing {
		t.Error("Should not be able to push items in a closing queue")
	}
	if _, err := queue.PushAsync("test", 10); err != ErrQueueIsClosing {
		t.Error("Should not be able to push items in a closing queue")
	}
}

func TestQueuePush(t *testing.T) {
	queue = &PriorityQueue{
		concurrency: 4,
		processorFn: process,
		in:          make(chan *Item, 1000),
		out:         make(chan map[uint32]interface{}, 100),
		quit:        make(chan bool),
		ready:       make(chan bool),
		wg:          &sync.WaitGroup{},
	}
	queue.Init()
	queue.WaitUntilReady()

	var err error

	data, err := queue.PushSync("testsync", 20)
	if err != nil {
		t.Error("There should be no error")
	}
	if data != "foo" {
		t.Error("Unexpected output received")
	}

	_, err = queue.PushAsync("testasync", 10)
	if err != nil {
		t.Error("There should be no error")
	}
}
