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
	"github.com/exago/svc/repository/processor"
)

var (
	ErrQueueIsClosing = errors.New("The queue is closing")

	queue *PriorityQueue
	once  sync.Once
)

type Worker struct {
	id   int
	in   chan *Item
	busy bool

	sync.RWMutex
}

type PriorityQueue struct {
	concurrency int
	items       ItemList
	workers     []*Worker
	in          chan *Item
	out         chan map[uint32]interface{}
	quit        chan bool
	closing     bool
	wg          *sync.WaitGroup
	processor   processor.Checker

	sync.RWMutex
}

func GetInstance() *PriorityQueue {
	log.SetLevel(log.DebugLevel)
	once.Do(func() {
		queue = &PriorityQueue{
			concurrency: 4,
			in:          make(chan *Item, 1000),
			out:         make(chan map[uint32]interface{}, 100),
			quit:        make(chan bool),
			wg:          &sync.WaitGroup{},
		}
		queue.Init()
	})
	return queue
}

func (pq *PriorityQueue) Init() {
	heap.Init(&pq.items)
	pq.InitWorkerPool()

	pq.wg.Add(1)
	go pq.Wait()

	go func() {
		sn := make(chan os.Signal, 1)
		signal.Notify(sn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sn
		pq.Stop()
	}()
}

func (pq *PriorityQueue) Wait() {
	defer pq.wg.Done()

	for {
		select {
		case in := <-pq.in:
			if in == nil {
				continue
			}
			log.Infoln("Received new item", in)
			heap.Push(&pq.items, in)
			pq.Process()
		case out := <-pq.out:
			log.Infoln("Item processed", out)
			pq.Process()
		case <-pq.quit:
			return
		}
	}
}

func (pq *PriorityQueue) Process() {
	pq.Lock()
	defer pq.Unlock()

	for pq.items.Len() > 0 {
		aw := pq.AvailableWorkers()
		if len(aw) == 0 {
			log.Debugln("Process", "No worker available")
			return
		}
		item := heap.Pop(&pq.items).(*Item)
		pq.PushToWorker(aw[0], item)
	}
}

func (pq *PriorityQueue) InitWorkerPool() {
	pq.wg.Add(pq.concurrency)
	for i := 0; i < pq.concurrency; i++ {
		go pq.Worker(i)
	}
}

func (pq *PriorityQueue) Worker(id int) {
	defer pq.wg.Done()
	w := Worker{id: id, in: make(chan *Item)}

	pq.Lock()
	pq.workers = append(pq.workers, &w)
	pq.Unlock()

	for {
		select {
		case item := <-w.in:
			log.Debugf("Worker %d got item %s with priority: %d", w.id, item.value, item.priority)

			// Fake processing
			time.Sleep(5 * time.Second)
			out := map[uint32]interface{}{}
			out[item.hash] = "foo"
			log.Debugf("Worker %d finished processing", w.id)

			pq.out <- out

			w.Lock()
			w.busy = false
			w.Unlock()
		case <-pq.quit:
			log.Infof("Worker %d got close", w.id)
			return
		}
	}
}

func (pq *PriorityQueue) AvailableWorkers() (workers []*Worker) {
	for _, worker := range pq.workers {
		worker.RLock()
		if worker.busy == false {
			workers = append(workers, worker)
		}
		worker.RUnlock()
	}
	log.Debugf("%d workers available", len(workers))
	return workers
}

func (pq *PriorityQueue) PushToWorker(worker *Worker, item *Item) {
	worker.Lock()
	defer worker.Unlock()
	log.Debugf("Pushed to worker %d: %.2d:%s", worker.id, item.priority, item.value)
	worker.busy = true
	worker.in <- item
}

func (pq *PriorityQueue) PushAsync(value string, priority int) (hash uint32, err error) {
	if pq.closing {
		return hash, ErrQueueIsClosing
	}
	item := Item{value: value, priority: priority}
	pq.in <- &item
	return item.Hash(), nil
}

func (pq *PriorityQueue) PushSync(value string, priority int) (hash uint32, err error) {
	if pq.closing {
		return hash, ErrQueueIsClosing
	}
	item := Item{value: value, priority: priority}
	pq.in <- &item
	return item.Hash(), nil
}

func (pq *PriorityQueue) Stop() {
	pq.closing = true
	close(pq.in)
	close(pq.quit)
	pq.wg.Wait()
	log.Debug("Queue safely stopped")
}
