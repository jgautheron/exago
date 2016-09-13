package repository

import (
	"encoding/json"
	"fmt"
)

// Save persists in database the repository data.
func (r Repository) Save() error {
	b, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	return r.DB.Put(r.cacheKey(), b)
}

// ClearCache removes the repository from database.
func (r Repository) ClearCache() error {
	return r.DB.Delete(r.cacheKey())
}

// IsCached checks if the repository's data is cached in database.
func (r Repository) IsCached() bool {
	if _, err := r.DB.Get(r.cacheKey()); err != nil {
		return false
	}
	return true
}

// cacheKey returns the standardised key format.
func (r Repository) cacheKey() []byte {
	return []byte(fmt.Sprintf("%s-%s", r.Name, r.Branch))
}
