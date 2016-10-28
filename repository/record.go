package repository

import (
	"time"

	"github.com/hotolab/exago-svc/repository/model"
)

type Record interface {
	GetName() string
	GetRank() string
	GetData() model.Data
	GetMetadata() model.Metadata
	SetMetadata() (err error)
	GetLastUpdate() time.Time
	SetLastUpdate(t time.Time)
	GetExecutionTime() string
	SetExecutionTime(duration time.Duration)
	GetScore() model.Score
	SetScore() (err error)
	GetCodeStats() model.CodeStats
	SetCodeStats(model.CodeStats)
	GetProjectRunner() model.ProjectRunner
	SetProjectRunner(tr model.ProjectRunner)
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
