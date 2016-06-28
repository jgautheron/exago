// Package showcaser processes repositories to determine occurrences, top-k, popularity.
package showcaser

import (
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/dgryski/go-topk"
	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/repository"
)

const (
	ItemCount   = 6
	TopkCount   = 100
	DatabaseKey = "showcase"
)

var (
	data    Showcase
	signals = make(chan os.Signal, 1)
)

type serialisedShowcase struct {
	Recent    []string
	TopRanked []string
	Popular   []string
	Topk      []byte
}

type Showcase struct {
	recent    []repository.RepositoryData
	topRanked []repository.RepositoryData
	popular   []repository.RepositoryData

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
func (d *Showcase) AddRecent(repo repository.RepositoryData) {
	d.Lock()
	defer d.Unlock()

	// Prevent duplicates
	for _, item := range d.recent {
		if item.GetName() == repo.GetName() {
			return
		}
	}

	d.recent = append(d.recent, repo)
	if len(d.recent) > d.itemCount {
		d.recent = d.recent[1:]
	}
}

// AddTopRanked pushes to the stack latest new A-ranked items, pops the old ones.
func (d *Showcase) AddTopRanked(repo repository.RepositoryData) {
	d.Lock()
	defer d.Unlock()

	if repo.GetRank() != "A" {
		return
	}

	// Prevent duplicates
	for _, item := range d.topRanked {
		if item.GetName() == repo.GetName() {
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
func (d *Showcase) AddPopular(repo repository.RepositoryData) {
	d.Lock()
	defer d.Unlock()

	d.tk.Insert(repo.GetName(), 1)
}

// updatePopular rebuilds the popular list from the topk stream.
// Doing it in real time would be unecessarily too expensive.
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
		sc.Recent = append(sc.Recent, r.GetName())
	}

	for _, r := range d.topRanked {
		sc.TopRanked = append(sc.TopRanked, r.GetName())
	}

	for _, r := range d.popular {
		sc.Popular = append(sc.Popular, r.GetName())
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
	return leveldb.Save([]byte(DatabaseKey), b)
}

func loadReposFromList(list []string) (repos []repository.RepositoryData, err error) {
	for _, name := range list {
		rp := repository.New(name, "")
		if err = rp.Load(); err != nil {
			return nil, err
		}
		repos = append(repos, rp)
	}
	return
}

// loadFromDB attempts to load a previously saved snapshot.
func loadFromDB() (s Showcase, exists bool, err error) {
	var repos []repository.RepositoryData

	b, err := leveldb.Get([]byte(DatabaseKey))
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

func ProcessRepository(repo repository.RepositoryData) {
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
