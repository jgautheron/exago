package repository

import (
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
}

func New(repo, branch string) model.Record {
	return &Repository{
		Name:   repo,
		Branch: branch,
	}
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
