package queue

import (
	"sync"
	"testing"
)

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
}
