// Package showcaser processes repositories to determine occurrences, top-k, popularity.
package showcaser

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/dgryski/go-topk"
	"github.com/hotolab/exago-svc/leveldb"
	"github.com/hotolab/exago-svc/repository"
	st "github.com/palantir/stacktrace"
	ldb "github.com/syndtr/goleveldb/leveldb"
)

const (
	ItemCount   = 6
	TopkCount   = 100
	DatabaseKey = "showcase"
)

var (
	showcase *Showcase
	once     sync.Once
	logger   = log.WithField("prefix", "showcaser")
	signals  = make(chan os.Signal, 1)
)

type serialisedShowcase struct {
	Recent    []string
	TopRanked []string
	Popular   []string
	Topk      []byte
}

type Showcase struct {
	recent    []repository.Record
	topRanked []repository.Record
	popular   []repository.Record

	// How many items per category
	itemCount int
	tk        *topk.Stream

	db leveldb.Database
	sync.RWMutex
}

// AddRecent pushes to the stack latest new items, pops the old ones.
func (d *Showcase) AddRecent(repo repository.Record) {
	d.Lock()
	defer d.Unlock()

	// Prevent duplicates
	for _, item := range d.recent {
		if item.GetName() == repo.GetName() {
			return
		}
	}

	d.recent = append(d.recent, repo)
	// Trim the list if it's too big
	if len(d.recent) > d.itemCount {
		d.recent = d.recent[1:]
	}
}

// AddTopRanked pushes to the stack latest new A-ranked items, pops the old ones.
func (d *Showcase) AddTopRanked(repo repository.Record) {
	d.Lock()
	defer d.Unlock()

	if !strings.Contains(repo.GetRank(), "A") {
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
func (d *Showcase) AddPopular(repo repository.Record) {
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
		if i < d.itemCount {
			list = append(list, v.Key)
		}
	}

	repos, err := d.loadReposFromList(list)
	if err != nil {
		return st.Propagate(err, "Failed to load the repos from the list: %#v", repos)
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
		return st.Propagate(err, "Could not serialize the index")
	}
	return d.db.Put([]byte(DatabaseKey), b)
}

func (d *Showcase) loadReposFromList(list []string) (repos []repository.Record, err error) {
	if len(list) == 0 {
		return
	}
	for _, name := range list {
		rp := &repository.Repository{Name: name, DB: d.db}
		if err = rp.Load(); err != nil {
			return nil, err
		}
		repos = append(repos, rp)
	}
	return
}

// loadFromDB attempts to load a previously saved snapshot.
func (d *Showcase) loadFromDB() (s Showcase, exists bool, err error) {
	var repos []repository.Record

	b, err := d.db.Get([]byte(DatabaseKey))
	if err == ldb.ErrNotFound {
		return s, false, nil
	}
	if b == nil || err != nil {
		return s, false, st.Propagate(err, "An error occurred while trying to load the snapshot", err)
	}

	var sc serialisedShowcase
	if err = json.Unmarshal(b, &sc); err != nil {
		return s, false, st.Propagate(err, "Failed to unmarshal the snapshot", err)
	}

	repos, err = d.loadReposFromList(sc.Recent)
	if err != nil {
		return s, false, st.Propagate(err, "Could not load recent repos", err)
	}
	s.recent = repos

	repos, err = d.loadReposFromList(sc.Popular)
	if err != nil {
		return s, false, st.Propagate(err, "Could not load popular repos", err)
	}
	s.popular = repos

	repos, err = d.loadReposFromList(sc.TopRanked)
	if err != nil {
		return s, false, st.Propagate(err, "Could not load top repos", err)
	}
	s.topRanked = repos

	s.tk = topk.New(TopkCount)
	s.tk.GobDecode(sc.Topk)

	return s, true, nil
}

func (d *Showcase) Process(repo repository.Record) {
	d.AddRecent(repo)
	d.AddTopRanked(repo)
	d.AddPopular(repo)
}

func GetInstance() *Showcase {
	once.Do(func() {
		if err := Init(); err != nil {
			logger.Fatal(err)
		}
	})
	return showcase
}

func Init() (err error) {
	db := leveldb.GetInstance()
	showcase = &Showcase{
		itemCount: ItemCount,
		tk:        topk.New(TopkCount),
		db:        db,
	}

	snapshot, exists, err := showcase.loadFromDB()
	if err != nil {
		logger.Errorf("An error occurred while loading the snapshot: %v", err)
	} else if exists {
		showcase = &snapshot
		showcase.db = db
		showcase.itemCount = ItemCount
		logger.Info("Snapshot loaded")
	}

	go showcase.catchInterrupt()
	go showcase.periodicallyRebuildPopularList()
	// go data.periodicallySave()
	return
}
