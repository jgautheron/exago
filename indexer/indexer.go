// Package indexer processes repositories to determine occurrences, top-k, popularity.
package indexer

import (
	"time"

	"github.com/dgryski/go-topk"
	"github.com/exago/svc/repository"
)

var data IndexedData

type IndexedData struct {
	Recent  []repository.Repository
	Top     []repository.Repository
	Popular []repository.Repository

	// How many items per category
	items int
	tk    *topk.Stream
}

// AddRecent pushes to the stack latest new items, pops the old ones.
func (d *IndexedData) AddRecent(repo repository.Repository) {
	d.Recent = append(d.Recent, repo)
	if len(d.Recent) > d.items {
		d.Recent = d.Recent[1:]
	}
}

// AddTop inserts the repository name into the topk data structure.
func (d *IndexedData) AddTop(repo repository.Repository) {
	d.tk.Insert(repo.Name, 1)
}

// updateTopK rebuilds the topk from the stream periodically.
func (d *IndexedData) updateTopK() {
	for {
		time.Sleep(5 * time.Minute)

		top := []repository.Repository{}
		for i, v := range d.tk.Keys() {
			if i > d.items {
				rp := repository.New(v.Key, "")
				rp.Load()
				top = append(top, *rp)
			}
		}
		d.Top = top
	}
}

func ProcessRepository(repo repository.Repository) {
	data.AddRecent(repo)
	data.AddTop(repo)
}

func init() {
	data = IndexedData{
		items: 6,
		tk:    topk.New(20),
	}

	go data.updateTopK()
}
