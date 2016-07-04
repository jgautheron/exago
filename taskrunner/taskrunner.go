package taskrunner

type TaskRunner interface {
	FetchCodeStats() (interface{}, error)
	FetchImports() (interface{}, error)
	FetchLintMessages(linters []string) (interface{}, error)
	FetchTestResults() (interface{}, error)
}
