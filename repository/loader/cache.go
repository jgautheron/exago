package loader

import (
	"encoding/json"
	"fmt"

	"github.com/hotolab/exago-svc/repository/model"
)

// Save persists in database the repository data.
func (l Loader) Save(repo model.Record) error {
	b, err := json.Marshal(repo)
	if err != nil {
		return err
	}
	return l.config.DB.Put(getCacheKey(repo.GetName(), repo.GetBranch(), repo.GetGoVersion()), b)
}

// ClearCache removes the repository from database.
func (l Loader) ClearCache(repo, branch, goversion string) error {
	return l.config.DB.Delete(getCacheKey(repo, branch, goversion))
}

// IsCached checks if the repository's data is cached in database.
func (l Loader) IsCached(repo, branch, goversion string) bool {
	if _, err := l.config.DB.Get(getCacheKey(repo, branch, goversion)); err != nil {
		return false
	}
	return true
}

// cacheKey returns the standardised key format.
func getCacheKey(repo, branch, goversion string) []byte {
	return []byte(fmt.Sprintf("%s|%s|%s", repo, branch, goversion))
}
