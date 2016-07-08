package repository

import (
	"encoding/json"
	"fmt"
)

// Save persists in database the repository data.
func (r *Repository) Save() error {
	b, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	return r.db.Put(r.cacheKey(), b)
}

// cacheKey returns the standardised key format.
func (r *Repository) cacheKey() []byte {
	return []byte(fmt.Sprintf("%s-%s", r.name, r.branch))
}
