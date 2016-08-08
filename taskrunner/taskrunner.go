package taskrunner

import "github.com/hotolab/exago-svc/repository/model"

type TaskRunner interface {
	FetchCodeStats() (model.CodeStats, error)
	FetchImports() (model.Imports, error)
	FetchLintMessages(linters []string) (model.LintMessages, error)
	FetchTestResults() (model.TestResults, error)
}
