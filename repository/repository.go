package repository

import (
	"encoding/json"
	"time"

	"github.com/go-xweb/log"
	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/score"
)

var (
	// Make sure it satisfies the interface.
	_ model.Record = (*Repository)(nil)
)

type Repository struct {
	Name, Branch string
	Data         model.Data
	startTime    time.Time
	loaded       bool

	DB model.Database
}

func New(repo, branch string, db model.Database, repositoryHost model.RepositoryHost) model.Record {
	return &Repository{
		Name:   repo,
		Branch: branch,
		DB:     db,
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
	r.Data.Name = r.Name
	r.Data.Branch = r.Branch
	r.loaded = true

	return nil
}

// ApplyScore calculates the score based on the repository results.
func (r *Repository) ApplyScore() (err error) {
	val, res := score.Process(r.Data)
	r.Data.Score.Value = val
	r.Data.Score.Details = res
	r.Data.Score.Rank = score.Rank(r.Data.Score.Value)

	log.Infof(
		"[%s] Rank: %s, overall score: %.2f",
		r.GetName(),
		r.Data.Score.Rank,
		r.Data.Score.Value,
	)

	return nil
}
