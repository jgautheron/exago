package repository

import (
	"fmt"
	"reflect"
	"time"

	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/repository/model"
)

type Repository struct {
	Name, Branch string

	// Data types
	CodeStats    model.CodeStats
	Imports      model.Imports
	TestResults  model.TestResults
	LintMessages model.LintMessages

	Score                 Score
	StartTime, LastUpdate time.Time
	ExecutionTime         string
}

func New(repo, branch string) *Repository {
	return &Repository{
		Name:   repo,
		Branch: branch,
	}
}

// IsCached checks if the repository's data is cached in database.
func (r *Repository) IsCached() bool {
	prefix := fmt.Sprintf("%s-%s", r.Name, r.Branch)
	data, err := leveldb.FindAllForRepository([]byte(prefix))
	if err != nil || len(data) != 7 {
		return false
	}
	return true
}

// IsLoaded checks if the data is already loaded.
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

// Load retrieves the entire matching dataset from database.
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
	if _, err = r.GetExecutionTime(); err != nil {
		return err
	}
	return err
}

// ClearCache removes the repository from database.
func (r *Repository) ClearCache() (err error) {
	prefix := fmt.Sprintf("%s-%s", r.Name, r.Branch)
	return leveldb.DeleteAllMatchingPrefix([]byte(prefix))
}

// FormatOutput prepares a map ready for output.
func (r *Repository) FormatOutput() map[string]interface{} {
	return map[string]interface{}{
		"imports":        r.Imports,
		"codestats":      r.CodeStats,
		"lintmessages":   r.LintMessages,
		"testresults":    r.TestResults,
		"score":          r.Score,
		"date":           r.LastUpdate,
		"execution_time": r.ExecutionTime,
	}
}
