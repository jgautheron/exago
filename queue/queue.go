// Package queue centralises all repositories processing.
package queue

import (
	"container/heap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

var (
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
	signal      chan os.Signal
	wg          *sync.WaitGroup

	sync.RWMutex
}

func GetInstance() *PriorityQueue {
	log.SetLevel(log.DebugLevel)
	once.Do(func() {
		queue = &PriorityQueue{
			concurrency: 4,
			in:          make(chan *Item),
			out:         make(chan map[uint32]interface{}),
			signal:      make(chan os.Signal),
			wg:          &sync.WaitGroup{},
		}
		queue.Init()
	})
	return queue
}

func (pq *PriorityQueue) Init() {
	signal.Notify(pq.signal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	heap.Init(&pq.items)
	pq.CreateWorkers()

	pq.wg.Add(1)
	go pq.Wait()
}

func (pq *PriorityQueue) Wait() {
	defer pq.wg.Done()
	defer pq.Stop()

	for {
		select {
		case in := <-pq.in:
			log.Infoln("Received new item", in)
			heap.Push(&pq.items, in)
			pq.Process()
		case out := <-pq.out:
			log.Infoln("Item processed", out)
			pq.Process()
		case <-pq.signal:
			pq.Stop()
			return
		}
	}
}

func (pq *PriorityQueue) Process() {
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

func (pq *PriorityQueue) CreateWorkers() {
	for i := 0; i < pq.concurrency; i++ {
		pq.wg.Add(1)
		go pq.Worker(i)
	}
}

func (pq *PriorityQueue) Worker(id int) {
	w := Worker{id: id, in: make(chan *Item)}
	pq.workers = append(pq.workers, &w)
	defer pq.wg.Done()
	log.Debugln("Worker created", id)

	for {
		select {
		case item := <-w.in:
			log.Infof("Worker %d got item %s", w.id, item.value)

			w.Lock()
			w.busy = true
			log.Debugf("Worker %d locked", w.id)
			w.Unlock()

			// Fake processing
			time.Sleep(5 * time.Second)
			out := map[uint32]interface{}{}
			out[item.hash] = "foo"

			log.Debugf("Worker %d finished processing", id)

			pq.out <- out

			w.Lock()
			w.busy = false
			w.Unlock()

			log.Debugf("Worker %d not busy anymore", id)
		case <-pq.signal:
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
	log.Debugf("Pushed to worker %d: %.2d:%s", worker.id, item.priority, item.value)
	worker.in <- item
}

func (pq *PriorityQueue) PushAsync(value string, priority int) (hash uint32) {
	item := Item{value: value, priority: priority}
	pq.in <- &item
	return item.Hash()
}

func (pq *PriorityQueue) Stop() {
	close(pq.in)
	close(pq.out)
	pq.wg.Wait()
}
