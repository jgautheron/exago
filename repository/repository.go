package repository

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/repository/model"
)

var (
	errMissingData = errors.New("Not enough data to calculate the rank")
	errNoDataFound = errors.New("No data found in database for this repository")
)

type Repository struct {
	Name, Branch string

	CodeStats    model.CodeStats
	Imports      model.Imports
	TestResults  model.TestResults
	LintMessages model.LintMessages

	Score Score

	LastUpdate time.Time
}

func New(repo, branch string) *Repository {
	return &Repository{
		Name:   repo,
		Branch: branch,
	}
}

func (r *Repository) IsCached() bool {
	prefix := fmt.Sprintf("%s-%s", r.Name, r.Branch)
	data, err := leveldb.FindAllForRepository([]byte(prefix))
	if err != nil || len(data) != 6 {
		return false
	}
	return true
}

func (r *Repository) IsLoaded() bool {
	if r.CodeStats == nil {
		return false
	}
	if r.Imports == nil {
		return false
	}
	if reflect.DeepEqual(r.TestResults, model.TestResults{}) {
		return false
	}
	if r.LintMessages == nil {
		return false
	}
	return true
}

func (r *Repository) Load() (err error) {
	if _, err = r.GetImports(); err != nil {
		return err
	}
	if _, err = r.GetCodeStats(); err != nil {
		return err
	}
	if _, err = r.GetLintMessages(DefaultLinters); err != nil {
		return err
	}
	if _, err = r.GetTestResults(); err != nil {
		return err
	}
	if _, err = r.GetScore(); err != nil {
		return err
	}
	if _, err = r.GetDate(); err != nil {
		return err
	}
	return err
}
