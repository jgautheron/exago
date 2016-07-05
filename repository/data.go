package repository

import (
	"time"

	"github.com/exago/svc/repository/model"
)

type RepositoryData interface {
	GetName() string
	GetMetadataDescription() string
	GetMetadataImage() string
	GetRank() string
	GetMetadata() (d model.Metadata, err error)
	SetMetadata() (err error)
	GetLastUpdate() (string, error)
	SetLastUpdate() (err error)
	GetExecutionTime() (string, error)
	SetExecutionTime() (err error)
	GetScore() (sc model.Score, err error)
	SetScore() (err error)
	GetImports() (model.Imports, error)
	SetImports(model.Imports)
	GetCodeStats() (model.CodeStats, error)
	SetCodeStats(model.CodeStats)
	GetTestResults() (tr model.TestResults, err error)
	SetTestResults(tr model.TestResults)
	GetLintMessages(linters []string) (model.LintMessages, error)
	SetLintMessages(model.LintMessages)
	SetStartTime(t time.Time)
	IsCached() bool
	IsLoaded() bool
	Load() (err error)
	ClearCache() (err error)
	AsMap() map[string]interface{}
}
