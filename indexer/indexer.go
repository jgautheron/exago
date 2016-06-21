// Package indexer processes repositories to determine occurrences, top-k, popularity.
package indexer

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgryski/go-topk"
	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/repository"
)

var data IndexedData

type IndexedData struct {
	recent    []repository.Repository
	topRanked []repository.Repository
	popular   []repository.Repository

	// How many items per category
	itemCount int
	tk        *topk.Stream
}

// AddRecent pushes to the stack latest new items, pops the old ones.
func (d *IndexedData) AddRecent(repo repository.Repository) {
	// Prevent duplicates
	for _, item := range d.recent {
		if item.Name == repo.Name {
			return
		}
	}

	d.recent = append(d.recent, repo)
	if len(d.recent) > d.itemCount {
		d.recent = d.recent[1:]
	}
}

// AddTopRanked pushes to the stack latest new A-ranked items, pops the old ones.
func (d *IndexedData) AddTopRanked(repo repository.Repository) {
	// Prevent duplicates
	for _, item := range d.topRanked {
		if item.Name == repo.Name {
			return
		}
	}

	d.topRanked = append(d.topRanked, repo)
	if len(d.topRanked) > 50 {
		d.topRanked = d.topRanked[1:]
	}
}

// AddPopular inserts the repository name into the topk data structure.
// The collection will not be updated in real-time (see updatePopular).
func (d *IndexedData) AddPopular(repo repository.Repository) {
	d.tk.Insert(repo.Name, 1)
}

// updatePopular rebuilds the data slice from the stream periodically.
func (d *IndexedData) updatePopular() {
	for {
		time.Sleep(5 * time.Minute)

		top := []repository.Repository{}
		for i, v := range d.tk.Keys() {
			if i <= d.itemCount {
				rp := repository.New(v.Key, "")
				rp.Load()
				top = append(top, *rp)
			}
		}
		d.popular = top
	}
}

// serialize the index as an easily loadable format.
func (d *IndexedData) serialize() ([]byte, error) {
	s := struct {
		recent    []string
		topRanked []string
		popular   []string
		tk        []byte
	}{}

	for _, r := range d.recent {
		s.recent = append(s.recent, r.Name)
	}

	for _, r := range d.topRanked {
		s.recent = append(s.topRanked, r.Name)
	}

	for _, r := range d.popular {
		s.recent = append(s.popular, r.Name)
	}

	tk, err := d.tk.GobEncode()
	if err != nil {
		return nil, err
	}

	return json.Marshal(tk)
}

// save persists periodically the index in database.
func (d *IndexedData) save() {
	for {
		time.Sleep(20 * time.Minute)

		b, err := d.serialize()
		if err != nil {
			log.Errorf("Error while serializing index: %v", err)
			return
		}

		leveldb.Save([]byte("indexer"), b)
		log.Debug("Index persisted in database")
	}
}

func ProcessRepository(repo repository.Repository) {
	data.AddRecent(repo)
	data.AddPopular(repo)
	data.AddTopRanked(repo)
}

func init() {
	data = IndexedData{
		itemCount: 6,
		tk:        topk.New(100),
	}

	go data.updatePopular()
	go data.save()
}
