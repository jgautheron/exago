package repository

import (
	"encoding/json"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/repository/model"
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
	name, branch string
	startTime    time.Time
	db           leveldb.Database
	loaded       bool
	Data         model.Data
}

func New(repo, branch string) *Repository {
	return &Repository{
		name:   repo,
		branch: branch,
		db:     leveldb.GetInstance(),
		Data: model.Data{
			Errors: make(map[string]error),
		},
	}
}

// IsLoaded checks if the data is already loaded.
func (r *Repository) IsLoaded() bool {
	return r.loaded
}

// Load retrieves the saved repository data from the database.
func (r *Repository) Load() error {
	b, err := r.db.Get(r.cacheKey())
	if err != nil {
		return err
	}

	var data model.Data
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	logrus.Infoln(r, data)
	r.Data = data
	r.loaded = true

	return nil
}
