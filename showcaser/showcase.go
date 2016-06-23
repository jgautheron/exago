// Package showcaser processes repositories to determine occurrences, top-k, popularity.
package showcaser

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgryski/go-topk"
	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/repository"
)

const (
	ItemCount = 6
	TopkCount = 100
)

var (
	data        Showcase
	databaseKey = []byte("indexer")
)

type serialisedShowcase struct {
	recent    []string
	topRanked []string
	popular   []string
	tk        []byte
}

type Showcase struct {
	recent    []repository.Repository
	topRanked []repository.Repository
	popular   []repository.Repository

	// How many items per category
	itemCount int
	tk        *topk.Stream
}

func New() Showcase {
	return Showcase{
		itemCount: ItemCount,
		tk:        topk.New(TopkCount),
	}
}

// AddRecent pushes to the stack latest new items, pops the old ones.
func (d *Showcase) AddRecent(repo repository.Repository) {
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
func (d *Showcase) AddTopRanked(repo repository.Repository) {
	// Prevent duplicates
	for _, item := range d.topRanked {
		if item.Name == repo.Name {
			return
		}
	}

	d.topRanked = append(d.topRanked, repo)
	if len(d.topRanked) > 100 {
		d.topRanked = d.topRanked[1:]
	}
}

// AddPopular inserts the repository name into the topk data structure.
// The collection will not be updated in real-time (see updatePopular).
func (d *Showcase) AddPopular(repo repository.Repository) {
	d.tk.Insert(repo.Name, 1)
}

// updatePopular rebuilds the data slice from the stream periodically.
func (d *Showcase) updatePopular() {
	for {
		time.Sleep(5 * time.Minute)

		list := make([]string, d.itemCount)
		for i, v := range d.tk.Keys() {
			if i <= d.itemCount {
				list = append(list, v.Key)
			}
		}

		repos, err := loadReposFromList(list)
		if err != nil {
			log.Errorf("An error occurred while loading the repos: %v", err)
			continue
		}

		d.popular = repos
	}
}

// serialize the index as an easily loadable format.
func (d *Showcase) serialize() ([]byte, error) {
	sc := serialisedShowcase{}

	for _, r := range d.recent {
		sc.recent = append(sc.recent, r.Name)
	}

	for _, r := range d.topRanked {
		sc.recent = append(sc.topRanked, r.Name)
	}

	for _, r := range d.popular {
		sc.recent = append(sc.popular, r.Name)
	}

	tk, err := d.tk.GobEncode()
	if err != nil {
		return nil, err
	}

	return json.Marshal(tk)
}

// save persists periodically the index in database.
func (d *Showcase) save() error {
	b, err := d.serialize()
	if err != nil {
		return err
	}

	return leveldb.Save(databaseKey, b)
}

func loadReposFromList(list []string) (repos []repository.Repository, err error) {
	for _, name := range list {
		rp := repository.New(name, "")
		if err = rp.Load(); err != nil {
			return nil, err
		}
		repos = append(repos, *rp)
	}
	return
}

// loadFromDB attempts to load a previously saved snapshot.
func loadFromDB() (s Showcase, exists bool, err error) {
	var repos []repository.Repository

	b, err := leveldb.Get(databaseKey)
	if b == nil || err != nil {
		return s, false, err
	}

	var sc serialisedShowcase
	if err = json.Unmarshal(b, &sc); err != nil {
		return s, false, err
	}

	repos, err = loadReposFromList(sc.recent)
	if err != nil {
		return s, false, err
	}
	s.recent = repos

	repos, err = loadReposFromList(sc.popular)
	if err != nil {
		return s, false, err
	}
	s.popular = repos

	repos, err = loadReposFromList(sc.topRanked)
	if err != nil {
		return s, false, err
	}
	s.topRanked = repos

	s.tk = topk.New(TopkCount)
	s.tk.GobDecode(sc.tk)

	return s, true, nil
}

// catchInterrupt traps termination signals to save a snapshot.
func catchInterrupt() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	select {
	case <-signals:
		log.Warn("Termination signal caught, saving the showcaser entries...")
		data.save()
	}
}

func periodicallySave() {
	for {
		time.Sleep(20 * time.Minute)

		if err := data.save(); err != nil {
			log.Errorf("Error while serializing index: %v", err)
			continue
		}

		log.Debug("Index persisted in database")
	}
}

func ProcessRepository(repo repository.Repository) {
	data.AddRecent(repo)
	data.AddTopRanked(repo)
	data.AddPopular(repo)
}

func Init() {
	data = New()
	snapshot, exists, err := loadFromDB()
	if err != nil {
		log.Errorf("An error occurred while loading the snapshot: %v", err)
	} else if exists {
		data = snapshot
		log.Info("Showcaser snapshot loaded")
	}

	go data.updatePopular()
	go periodicallySave()
	// go catchInterrupt()
}
