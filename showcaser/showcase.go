// Package showcaser processes repositories to determine occurrences, top-k, popularity.
package showcaser

import (
	"encoding/json"
	"os"
	"os/signal"
	"sync"
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
	databaseKey = []byte("showcase")
	signals     = make(chan os.Signal, 1)
)

type serialisedShowcase struct {
	Recent    []string
	TopRanked []string
	Popular   []string
	Topk      []byte
}

type Showcase struct {
	recent    []repository.Repository
	topRanked []repository.Repository
	popular   []repository.Repository

	// How many items per category
	itemCount int
	tk        *topk.Stream

	sync.RWMutex
}

func New() Showcase {
	return Showcase{
		itemCount: ItemCount,
		tk:        topk.New(TopkCount),
	}
}

// AddRecent pushes to the stack latest new items, pops the old ones.
func (d *Showcase) AddRecent(repo repository.Repository) {
	d.Lock()
	defer d.Unlock()

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
	d.Lock()
	defer d.Unlock()

	if repo.Score.Rank != "A" {
		return
	}

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
	d.Lock()
	defer d.Unlock()

	d.tk.Insert(repo.Name, 1)
}

// updatePopular rebuilds the popular list from the topk stream.
func (d *Showcase) updatePopular() error {
	d.Lock()
	defer d.Unlock()

	list := []string{}
	for i, v := range d.tk.Keys() {
		if i <= d.itemCount {
			list = append(list, v.Key)
		}
	}

	repos, err := loadReposFromList(list)
	if err != nil {
		return err
	}

	d.popular = repos
	return nil
}

// serialize the index as an easily loadable format.
func (d *Showcase) serialize() ([]byte, error) {
	sc := serialisedShowcase{}

	for _, r := range d.recent {
		sc.Recent = append(sc.Recent, r.Name)
	}

	for _, r := range d.topRanked {
		sc.TopRanked = append(sc.TopRanked, r.Name)
	}

	for _, r := range d.popular {
		sc.Popular = append(sc.Popular, r.Name)
	}

	tk, err := d.tk.GobEncode()
	if err != nil {
		return nil, err
	}
	sc.Topk = tk

	return json.Marshal(sc)
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

	repos, err = loadReposFromList(sc.Recent)
	if err != nil {
		return s, false, err
	}
	s.recent = repos

	repos, err = loadReposFromList(sc.Popular)
	if err != nil {
		return s, false, err
	}
	s.popular = repos

	repos, err = loadReposFromList(sc.TopRanked)
	if err != nil {
		return s, false, err
	}
	s.topRanked = repos

	s.tk = topk.New(TopkCount)
	s.tk.GobDecode(sc.Topk)

	return s, true, nil
}

// catchInterrupt traps termination signals to save a snapshot.
func catchInterrupt() {
	select {
	case <-signals:
		log.Warn("Termination signal caught, saving the showcaser entries...")
		err := data.save()
		if err != nil {
			log.Errorf("Got error while saving: %v", err)
		}
		close(signals)
	}
}

func periodicallyRebuildPopularList() {
	for {
		select {
		case <-signals:
			return
		case <-time.After(10 * time.Minute):
			err := data.updatePopular()
			if err != nil {
				log.Errorf("Got error while updating the popular list: %v", err)
			}
			log.Debug("Rebuilt the popular list")
		}
	}
}

func periodicallySave() {
	for {
		select {
		case <-signals:
			return
		case <-time.After(30 * time.Minute):
			if err := data.save(); err != nil {
				log.Errorf("Error while serializing index: %v", err)
				continue
			}
			log.Debug("Index persisted in database")
		}
	}
}

func ProcessRepository(repo repository.Repository) {
	data.AddRecent(repo)
	data.AddTopRanked(repo)
	data.AddPopular(repo)
}

func Init() (err error) {
	data = New()
	snapshot, exists, err := loadFromDB()
	if err != nil {
		log.Errorf("An error occurred while loading the snapshot: %v", err)
	} else if exists {
		data = snapshot
		log.Info("Showcaser snapshot loaded")
	}

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go catchInterrupt()
	go periodicallyRebuildPopularList()
	// go periodicallySave()
	return nil
}
