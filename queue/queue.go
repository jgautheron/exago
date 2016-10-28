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
	processorFn func(repo, branch string, tr taskrunner.TaskRunner) (interface{}, error)
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
			in:          make(chan string, 1000),
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
	go pq.ProcessAll()

	// Trap interruption signals
	go pq.StopOnSignal()
}

func (pq *Queue) StopOnSignal() {
	sn := make(chan os.Signal, 1)
	signal.Notify(sn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sn
	pq.Stop()
}

// Wait processes incoming queue items.
func (pq *Queue) Wait() {
	defer pq.wg.Done()
	for {
		select {
		case out := <-pq.out:
			for repository := range out {
				logger.WithField("repository", repository).Debug("Item processed")
			}
		case <-pq.quit:
			return
		default:
		}
	}
}

func (pq *Queue) ProcessAll() {
	defer pq.wg.Done()
	for {
		log.Warn("foo")
		pq.sem <- true
		in := <-pq.in
		pq.Process(in)
	}
}

// Process waits for an available slot, processes the item and frees the slot.
func (pq *Queue) Process(repo string) {
	defer pq.wg.Done()
	defer func() { <-pq.sem }()

	lgr := logger.WithFields(log.Fields{
		"repository": repo,
	})

	lgr.Debug("Begin processing")
	data, err := pq.processorFn(repo, "", lambda.Runner{Repository: repo})
	if err != nil {
		lgr.Error(err)
		return
	}
	out := map[string]interface{}{}
	out[repo] = data
	pq.out <- out
}

// WaitUntilEmpty blocks until the queue is empty.
func (pq *Queue) WaitUntilEmpty() {
	doneCh := make(chan bool)
	go func() {
		for {
			if len(pq.sem) == 0 {
				doneCh <- true
				return
			}
		}
	}()
	<-doneCh
}

// PushAsync adds an item to the queue asynchronously.
func (pq *Queue) PushAsync(repository string) error {
	if pq.closing {
		return ErrQueueIsClosing
	}
	pq.in <- repository

	// Read from the channel
	// go func() {
	// 	out := <-pq.out
	// 	for repository := range out {
	// 		logger.WithField("repository", repository).Debug("Item processed")
	// 	}
	// }()

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
			default:
			}
		}
	}()

	wg.Wait()
	return data, nil
}

// Stop gracefully closes the queue and all its workers.
func (pq *Queue) Stop() {
	logger.Debug("Stop called")
	pq.closing = true
	close(pq.in)
	close(pq.quit)
	pq.wg.Wait()
	logger.Debug("Queue safely stopped")
}
