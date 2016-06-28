package repository

import (
	"fmt"
	"reflect"
	"time"

	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/repository/model"
	"github.com/exago/svc/repository/score"
)

var (
	// DefaultLinters run by default in Lambda
	DefaultLinters = []string{
		"deadcode",
		"dupl",
		"errcheck",
		"goconst",
		"gocyclo",
		"gofmt",
		"goimports",
		"golint",
		"gosimple",
		"ineffassign",
		"staticcheck",
		"varcheck",
		"vet",
		"vetshadow",
	}
)

type RepositoryData interface {
	GetName() string
	GetMetadataDescription() string
	GetMetadataImage() string
	GetRank() string
	GetMetadata() (d model.Metadata, err error)
	GetLastUpdate() (string, error)
	GetExecutionTime() (string, error)
	GetScore() (sc model.Score, err error)
	GetImports() (model.Imports, error)
	GetCodeStats() (model.CodeStats, error)
	GetTestResults() (tr model.TestResults, err error)
	GetLintMessages(linters []string) (model.LintMessages, error)
	SetStartTime(t time.Time)
	IsCached() bool
	IsLoaded() bool
	Load() (err error)
	ClearCache() (err error)
	AsMap() map[string]interface{}
}

type Repository struct {
	Name, Branch string

	// Data types
	CodeStats    model.CodeStats
	Imports      model.Imports
	TestResults  model.TestResults
	LintMessages model.LintMessages
	Metadata     model.Metadata
	Score        model.Score

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
	if err != nil || len(data) != 8 {
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
	if _, err = r.GetMetadata(); err != nil {
		return err
	}
	if _, err = r.GetLastUpdate(); err != nil {
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

// AsMap generates a map out of repository fields.
func (r *Repository) AsMap() map[string]interface{} {
	return map[string]interface{}{
		model.ImportsName:       r.Imports,
		model.CodeStatsName:     r.CodeStats,
		model.LintMessagesName:  r.LintMessages,
		model.TestResultsName:   r.TestResults,
		model.ScoreName:         r.Score,
		model.MetadataName:      r.Metadata,
		model.LastUpdateName:    r.LastUpdate,
		model.ExecutionTimeName: r.ExecutionTime,
	}
}

func (r *Repository) GetName() string {
	return r.Name
}

func (r *Repository) GetMetadataDescription() string {
	return r.Metadata.Description
}

func (r *Repository) GetMetadataImage() string {
	return r.Metadata.Image
}

func (r *Repository) GetRank() string {
	return r.Metadata.Image
}

func (r *Repository) SetStartTime(t time.Time) {
	r.StartTime = t
}

func (r *Repository) calcScore() {
	val, res := score.Process(r.AsMap())
	r.Score.Value = val
	r.Score.Details = res
	r.Score.Rank = score.Rank(r.Score.Value)
}
