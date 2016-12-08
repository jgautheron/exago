package server

import (
	"strings"
	"sync"
)

type processingList struct {
	repos map[string]bool
	sync.Mutex
}

func (pl *processingList) add(key ...string) {
	repo := strings.Join(key, "-")
	pl.Lock()
	pl.repos[repo] = true
	pl.Unlock()
}

func (pl *processingList) exists(key ...string) bool {
	repo := strings.Join(key, "-")
	_, found := pl.repos[repo]
	return found
}

func (pl *processingList) remove(key ...string) {
	repo := strings.Join(key, "-")
	pl.Lock()
	delete(pl.repos, repo)
	pl.Unlock()
}
