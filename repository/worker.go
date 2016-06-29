package repository

type Worker interface {
	GetCodeStats(repository string) (interface{}, error)
	GetImports(repository string) (interface{}, error)
	GetLintMessages(repository string, linters []string) (interface{}, error)
	GetTestResults(repository string) (interface{}, error)
}
