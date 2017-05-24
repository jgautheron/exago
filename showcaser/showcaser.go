// Package showcaser processes repositories to determine occurrences, top-k, popularity.
package showcaser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/dgryski/go-topk"
	exago "github.com/jgautheron/exago"
	"github.com/jgautheron/exago/repository/model"
	st "github.com/palantir/stacktrace"
	ldb "github.com/syndtr/goleveldb/leveldb"
)

const (
	ItemCount   = 6
	TopkCount   = 100
	DatabaseKey = "showcaser"
)

var (
	logger  = log.WithField("prefix", "showcaser")
	signals = make(chan os.Signal, 1)

	// Make sure it satisfies the interface.
	_ model.Promoter = (*Showcaser)(nil)
)

type Showcaser struct {
	recent    []model.Record
	topRanked []model.Record
	popular   []model.Record

	// How many items per category
	itemCount int
	tk        *topk.Stream

	config exago.Config
	sync.RWMutex
}

func New(options ...exago.Option) (*Showcaser, error) {
	showcaser := &Showcaser{
		itemCount: ItemCount,
		tk:        topk.New(TopkCount),
	}

	// Apply options
	for _, option := range options {
		option.Apply(&showcaser.config)
	}

	snapshot, exists, err := showcaser.loadFromDB()
	if err != nil {
		return nil, fmt.Errorf("An error occurred while loading the snapshot: %v", err)
	} else if exists {
		s := &snapshot
		s.itemCount = ItemCount
		s.config = showcaser.config
		logger.Info("Snapshot loaded")
		return s, nil
	}

	return showcaser, nil
}

func (d *Showcaser) StartRoutines() {
	go d.catchInterrupt()
	go d.periodicallyRebuildPopularList()
}

// AddRecent pushes to the stack latest new items, pops the old ones.
func (d *Showcaser) AddRecent(repo model.Record) {
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
func (d *Showcaser) AddTopRanked(repo model.Record) {
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
func (d *Showcaser) AddPopular(repo model.Record) {
	d.Lock()
	d.tk.Insert(repo.GetName(), 1)
	d.Unlock()
}

// updatePopular rebuilds the popular list from the topk stream.
// Doing it in real time would be unecessarily too expensive.
func (d *Showcaser) updatePopular() error {
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
func (d *Showcaser) serialize() ([]byte, error) {
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
func (d *Showcaser) save() error {
	b, err := d.serialize()
	if err != nil {
		return st.Propagate(err, "Could not serialize the index")
	}
	return d.config.DB.Put([]byte(DatabaseKey), b)
}

func (d *Showcaser) loadReposFromList(list []string) (repos []model.Record, err error) {
	if len(list) == 0 {
		return
	}
	for _, name := range list {
		rp, err := d.config.RepositoryLoader.Load(name, "")
		if err != nil {
			return nil, err
		}
		repos = append(repos, rp)
	}
	return
}

// loadFromDB attempts to load a previously saved snapshot.
func (d *Showcaser) loadFromDB() (s Showcaser, exists bool, err error) {
	var repos []model.Record

	b, err := d.config.DB.Get([]byte(DatabaseKey))
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

func (d *Showcaser) Process(repo model.Record) {
	d.AddRecent(repo)
	d.AddTopRanked(repo)
	d.AddPopular(repo)
}

type serialisedShowcase struct {
	Recent    []string
	TopRanked []string
	Popular   []string
	Topk      []byte
}
