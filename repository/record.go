package repository

import (
	"time"

	"github.com/exago/svc/repository/model"
)

type Record interface {
	GetName() string
	GetRank() string
	GetData() model.Data
	GetMetadata() model.Metadata
	SetMetadata() (err error)
	GetLastUpdate() time.Time
	SetLastUpdate()
	GetExecutionTime() string
	SetExecutionTime()
	GetScore() model.Score
	SetScore() (err error)
	GetImports() model.Imports
	SetImports(model.Imports)
	GetCodeStats() model.CodeStats
	SetCodeStats(model.CodeStats)
	GetTestResults() model.TestResults
	SetTestResults(tr model.TestResults)
	GetLintMessages() model.LintMessages
	SetLintMessages(model.LintMessages)
	SetStartTime(t time.Time)
	SetError(tp string, err error)
	IsCached() bool
	IsLoaded() bool
	Load() (err error)
	ClearCache() (err error)
	Save() error
}
