// Package gosearch is used for indexing all repositories.
// http://go-search.org/
package gosearch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	exago "github.com/hotolab/exago-svc"
)

const (
	RequestTimeout = time.Duration(10 * time.Second)
	URL            = "http://go-search.org/api?action=packages"
	IndexDBKey     = "gosearch:index"
)

var (
	logger = log.WithField("prefix", "gosearch")
)

type Gosearch struct {
	config exago.Config
}

func New(options ...exago.Option) *Gosearch {
	var gs Gosearch
	for _, option := range options {
		option.Apply(&gs.config)
	}
	return &gs
}

// LoadIndex queries go-search.org API and persists a list of unique GitHub repositories.
func (g *Gosearch) LoadIndex() error {
	c := http.Client{
		Timeout: RequestTimeout,
	}
	resp, err := c.Get(URL)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var packages []string
	if err := decoder.Decode(&packages); err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the list and automatically trim down duplicates with a map
	reposMap := map[string]bool{}
	r, _ := regexp.Compile(`^github.com/([\w\d\-]+)/([\w\d\-]+)`)
	for _, pkg := range packages {
		m := r.FindStringSubmatch(pkg)
		if len(m) == 0 {
			continue
		}
		reposMap[m[0]] = true
	}

	logger.Infof("Found %d unique GitHub repositories", len(reposMap))

	repos := []string{}
	for repo := range reposMap {
		repos = append(repos, repo)
	}

	// Persist the list in database as JSON
	b, err := json.Marshal(repos)
	if err != nil {
		return err
	}
	if err := g.config.DB.Put([]byte(IndexDBKey), b); err != nil {
		return fmt.Errorf("An error occurred while saving the index in DB: %v", err)
	}

	logger.Info("Successfully persisted the index")
	return nil
}

// GetIndex retrieves the GoSearch Index from the database.
func (g *Gosearch) GetIndex() (repos []string, err error) {
	b, err := g.config.DB.Get([]byte(IndexDBKey))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &repos)
	return repos, err
}
