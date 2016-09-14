// Package queue centralises all repositories processing.
package queue

import (
	"container/heap"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/repository/processor"
	"github.com/hotolab/exago-svc/taskrunner"
	"github.com/hotolab/exago-svc/taskrunner/lambda"
)

var (
	ErrQueueIsClosing = errors.New("The queue is closing")
	logger            = log.WithField("prefix", "queue")

	queue *PriorityQueue
	once  sync.Once
)

// Worker is a single processing unit, running forever.
// It simply waits for messages and processes them.
type Worker struct {
	id   int
	in   chan *Item
	busy bool

	sync.RWMutex
}

// PriorityQueue is a queue based on a heap.
// Its messages can be distributed to many workers (defined by concurrency).
type PriorityQueue struct {
	concurrency int
	items       ItemList
	workers     []*Worker
	processorFn func(value string, tr taskrunner.TaskRunner) (interface{}, error)
	in          chan *Item
	out         chan map[uint32]interface{}
	quit        chan bool
	ready       chan bool
	closing     bool
	wg          *sync.WaitGroup

	sync.RWMutex
}

// GetInstance returns the queue instance.
// The queue will be instantiated if it wasn't yet.
func GetInstance() *PriorityQueue {
	once.Do(func() {
		queue = &PriorityQueue{
			concurrency: 4,
			processorFn: processor.ProcessRepository,
			in:          make(chan *Item, 1000),
			out:         make(chan map[uint32]interface{}, 100),
			quit:        make(chan bool),
			wg:          &sync.WaitGroup{},
		}
		queue.Init()
	})
	return queue
}

// Init lays the ground work necessary for the queue to function properly.
func (pq *PriorityQueue) Init() {
	heap.Init(&pq.items)
	pq.InitWorkerPool()

	go func() {
		for {
			if len(pq.workers) == pq.concurrency {
				logger.Debug("Ready")
				pq.ready <- true
			}
		}
	}()

	pq.wg.Add(1)
	go pq.Wait()

	// Trap interruption signals
	go func() {
		sn := make(chan os.Signal, 1)
		signal.Notify(sn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sn
		pq.Stop()
	}()
}

// Wait listens for new and processed items.
// - New items are pushed in the heap and will be eventually processed.
// - Asynchronously processed items are received here, no further action needed.
func (pq *PriorityQueue) Wait() {
	defer pq.wg.Done()

	for {
		select {
		case in := <-pq.in:
			if in == nil {
				continue
			}
			logger.WithField("hash", in.hash).Debug("Received new item")
			heap.Push(&pq.items, in)
			pq.Process()
		case out := <-pq.out:
			for hash, _ := range out {
				logger.WithField("hash", hash).Debug("Item processed")
			}
			// A worker just got freed, give it the next available item
			pq.Process()
		case <-pq.quit:
			return
		}
	}
}

// Process attempts to push the next item to the first available worker.
func (pq *PriorityQueue) Process() {
	pq.Lock()
	defer pq.Unlock()
	for pq.items.Len() > 0 {
		aw := pq.AvailableWorkers()
		if len(aw) == 0 {
			logger.Debug("No worker available")
			return
		}
		item := heap.Pop(&pq.items).(*Item)
		pq.PushToWorker(aw[0], item)
	}
}

// InitWorkerPool creates the initial worker pool, the amount of workers being
// defined by the concurrency setting.
// It cannot exceed 25 (4 per check) since AWS Lambda has a default safety
// throttle set to 100 concurrent executions.
func (pq *PriorityQueue) InitWorkerPool() {
	pq.wg.Add(pq.concurrency)
	for i := 0; i < pq.concurrency; i++ {
		go pq.Worker(i)
	}
}

// Worker runs a queue worker in background.
// Its sole purpose is to process messages, once a message is processed
// the output it caught by Wait().
func (pq *PriorityQueue) Worker(id int) {
	defer pq.wg.Done()
	w := Worker{id: id, in: make(chan *Item)}

	pq.Lock()
	pq.workers = append(pq.workers, &w)
	pq.Unlock()

	for {
		select {
		case item := <-w.in:
			data, _ := pq.processorFn(item.value, lambda.Runner{Repository: item.value})

			out := map[uint32]interface{}{}
			out[item.hash] = data
			pq.out <- out

			logger.WithField("hash", item.hash).Debugf("Worker %d finished processing", w.id)
			w.Lock()
			w.busy = false
			w.Unlock()
		case <-pq.quit:
			logger.Debugf("Stopping worker %d", w.id)
			pq.removeWorker(w)
			return
		}
	}
}

// AvailableWorkers returns the workers which are currently available
// for processing.
func (pq *PriorityQueue) AvailableWorkers() (workers []*Worker) {
	for _, worker := range pq.workers {
		worker.RLock()
		if worker.busy == false {
			workers = append(workers, worker)
		}
		worker.RUnlock()
	}
	logger.Debugf("%d workers available", len(workers))
	return workers
}

// PushToWorker assigns an item to be processed to the given worker.
func (pq *PriorityQueue) PushToWorker(worker *Worker, item *Item) {
	worker.Lock()
	defer worker.Unlock()
	logger.WithField("hash", item.hash).Debugf(
		"Pushed to worker %d",
		worker.id,
	)
	worker.busy = true
	worker.in <- item
}

// WaitUntilReady blocks until the queue is properly initialized.
func (pq *PriorityQueue) WaitUntilReady() {
	select {
	case <-pq.ready:
		return
	case <-time.After(1 * time.Second):
		logger.Error("The queue took too long to initiate (1s)")
	}
}

// PushAsync adds an item to the queue asynchronously.
func (pq *PriorityQueue) PushAsync(value string, priority int) (hash uint32, err error) {
	if pq.closing {
		return hash, ErrQueueIsClosing
	}
	item := NewItem(value, priority)
	pq.in <- item
	return item.hash, nil
}

// PushSync adds an item to the queue synchronously, blocking until
// the processing is done.
func (pq *PriorityQueue) PushSync(value string, priority int) (data interface{}, err error) {
	if pq.closing {
		return nil, ErrQueueIsClosing
	}
	item := NewItem(value, priority)
	pq.in <- item

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case out := <-pq.out:
				for hash, out := range out {
					if hash == item.hash {
						data = out
						return
					}
				}
			case <-pq.quit:
				return
			}
		}
	}()

	wg.Wait()
	return data, nil
}

// Stop gracefully closes the queue and all its workers.
func (pq *PriorityQueue) Stop() {
	pq.closing = true
	close(pq.in)
	close(pq.quit)
	pq.wg.Wait()
	logger.Debug("Queue safely stopped")
}

// removeWorker removes the given worker from the list of available workers.
func (pq *PriorityQueue) removeWorker(worker Worker) {
	pq.Lock()
	defer pq.Unlock()

	workers := []*Worker{}
	for _, w := range pq.workers {
		if w.id != worker.id {
			workers = append(workers, w)
		}
	}
	pq.workers = workers
}
