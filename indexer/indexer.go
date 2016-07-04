// Package indexer enables mass processing of repositories.
package indexer

import (
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/repository/processor"
	"github.com/exago/svc/taskrunner/lambda"
)

type Indexer struct {
	Done, Aborted chan bool

	items           map[string]itemState
	processingItems map[string]bool
	concurrency     int
	processedItems  chan processedItem
	processedCount  int

	sync.RWMutex
}

// New prepares a new indexer.
func New(items []string) *Indexer {
	c := runtime.NumCPU()

	mp := map[string]itemState{}
	for _, item := range items {
		mp[item] = itemState{nil, false}
	}

	return &Indexer{
		Done:            make(chan bool, 1),
		Aborted:         make(chan bool, 1),
		items:           mp,
		processingItems: make(map[string]bool, c),
		concurrency:     c,
		processedItems:  make(chan processedItem, len(mp)),
		processedCount:  0,
	}
}

// Start runs the indexer.
func (idx *Indexer) Start() {
	start := time.Now()

	go func() {
		for item := range idx.processedItems {
			state := itemState{nil, true}
			if item.err != nil {
				state = itemState{item.err, true}
			}

			idx.items[item.name] = state
			delete(idx.processingItems, item.name)
			idx.processedCount++

			// If all items are processed, we're done
			if idx.processedCount == len(idx.items) {
				idx.Done <- true
				break
			}

			// Process the next available item
			go idx.ProcessItem()
		}
	}()

	for i := 0; i < idx.concurrency; i++ {
		go idx.ProcessItem()
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	select {
	case <-idx.Done:
		elapsed := time.Since(start)
		log.Infof("Processed %d item(s) in %s", len(idx.items), elapsed)
	case <-signals:
		log.Warn("Termination signal caught, stopping the indexer")
		idx.Aborted <- true
	}
}

// ProcessItem evaluates the next available item.
func (idx *Indexer) ProcessItem() {
	for item, state := range idx.items {
		// Has the item already been processed?
		if state.processed {
			continue
		}

		// Is the item currently being processed?
		idx.RLock()
		if _, beingProcessed := idx.processingItems[item]; beingProcessed {
			continue
		}
		idx.RUnlock()
		idx.processingItems[item] = true

		lgr := log.WithField("repository", item)
		lgr.Info("Processing...")

		rc := processor.NewChecker(item, lambda.Runner{})
		if rc.Repository.IsCached() {
			// Possibly useful later: a --force flag to reprocess already cached repos
			lgr.Warn("Already processed")
			idx.processedItems <- processedItem{item, nil}
			continue
		}

		go func() {
			// If an error is caught, abort the processing
			for err := range rc.Errors {
				state.err = err
				rc.Abort()
				break
			}
		}()

		// Process the repository
		go rc.Run()

		select {
		case <-rc.Aborted:
			lgr.WithField("error", state.err).Warn("Processing aborted")
			idx.processedItems <- processedItem{item, state.err}
		case <-rc.Done:
			lgr.WithField("score", rc.Repository.GetRank()).Info("Processing successful")
			idx.processedItems <- processedItem{item, nil}
		case <-idx.Aborted:
			rc.Abort()
		}
	}
}

type itemState struct {
	err       error
	processed bool
}

type processedItem struct {
	name string
	err  error
}
