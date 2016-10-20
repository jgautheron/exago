package taskrunner

import "github.com/hotolab/exago-svc/repository/model"

type TaskRunner interface {
	FetchCodeStats() (model.CodeStats, error)
	FetchLintMessages() (model.LintMessages, error)
	FetchProjectRunner() (model.ProjectRunner, error)
}
