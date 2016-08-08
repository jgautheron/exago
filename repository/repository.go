package repository

import (
	"encoding/json"
	"time"

	"github.com/hotolab/exago-svc/github"
	"github.com/hotolab/exago-svc/leveldb"
	"github.com/hotolab/exago-svc/repository/model"
)

var (
	// DefaultLinters ran by default in Lambda.
	DefaultLinters = []string{
		"deadcode", "dupl", "errcheck", "goconst", "gocyclo", "gofmt", "goimports",
		"golint", "gosimple", "ineffassign", "staticcheck", "vet", "vetshadow",
	}

	// Make sure it satisfies the interface.
	_ Record = (*Repository)(nil)
)

type Repository struct {
	Name, Branch   string
	DB             leveldb.Database
	RepositoryHost github.RepositoryHost
	Data           model.Data
	startTime      time.Time
	loaded         bool
}

func New(repo, branch string) *Repository {
	return &Repository{
		Name:           repo,
		Branch:         branch,
		DB:             leveldb.GetInstance(),
		RepositoryHost: github.GetInstance(),
	}
}

// IsLoaded checks if the data is already loaded.
func (r *Repository) IsLoaded() bool {
	return r.loaded
}

// Load retrieves the saved repository data from the database.
func (r *Repository) Load() error {
	b, err := r.DB.Get(r.cacheKey())
	if err != nil {
		return err
	}

	var data model.Data
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	r.Data = data
	r.loaded = true

	return nil
}
