package queue

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/repository/processor"
	"github.com/hotolab/exago-svc/taskrunner"
	"github.com/hotolab/exago-svc/taskrunner/lambda"
)

var (
	ErrQueueIsClosing = errors.New("The queue is closing")
	logger            = log.WithField("prefix", "queue")

	queue *Queue
	once  sync.Once
)

type Queue struct {
	processorFn func(value string, tr taskrunner.TaskRunner) (interface{}, error)
	in          chan string
	out         chan map[string]interface{}
	quit        chan bool
	sem         chan bool
	closing     bool
	wg          *sync.WaitGroup
}

// GetInstance returns the queue instance.
// The queue will be instantiated if it wasn't yet.
func GetInstance() *Queue {
	once.Do(func() {
		queue = &Queue{
			processorFn: processor.ProcessRepository,
			sem:         make(chan bool, 4),
			in:          make(chan string),
			out:         make(chan map[string]interface{}),
			quit:        make(chan bool),
			wg:          &sync.WaitGroup{},
		}
		queue.Init()
	})
	return queue
}

// Init lays the ground work.
func (pq *Queue) Init() {
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

// Wait processes incoming queue items.
func (pq *Queue) Wait() {
	defer pq.wg.Done()
	for {
		select {
		case in := <-pq.in:
			if in == "" {
				continue
			}
			pq.wg.Add(1)
			go pq.Process(in)
		case out := <-pq.out:
			for repository := range out {
				logger.WithField("repository", repository).Debug("Item processed")
			}
		case <-pq.quit:
			return
		}
	}
}

// Process waits for an available slot, processes the item and frees the slot.
func (pq *Queue) Process(repository string) {
	pq.sem <- true
	defer pq.wg.Done()
	defer func() { <-pq.sem }()

	logger.WithFields(log.Fields{
		"repository": repository,
	}).Debug("Beginning processing")

	data, err := pq.processorFn(repository, lambda.Runner{Repository: repository})
	if err != nil {
		logger.Error(err)
		return
	}
	out := map[string]interface{}{}
	out[repository] = data
	pq.out <- out
}

// WaitUntilEmpty blocks until the queue is empty.
func (pq *Queue) WaitUntilEmpty() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			if len(pq.sem) == 0 {
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
}

// PushAsync adds an item to the queue asynchronously.
func (pq *Queue) PushAsync(repository string) error {
	if pq.closing {
		return ErrQueueIsClosing
	}
	pq.in <- repository
	return nil
}

// PushSync adds an item to the queue synchronously, blocking until
// the processing is done.
func (pq *Queue) PushSync(repository string) (data interface{}, err error) {
	if pq.closing {
		return nil, ErrQueueIsClosing
	}
	pq.in <- repository

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case out := <-pq.out:
				for rp, out := range out {
					if rp == repository {
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
func (pq *Queue) Stop() {
	pq.closing = true
	close(pq.in)
	close(pq.quit)
	pq.wg.Wait()
	logger.Debug("Queue safely stopped")
}
