package datafetcher

import (
	"encoding/json"
	"strings"
)

func GetLintResults(repository, linter string) (*json.RawMessage, error) {
	sp := strings.Split(repository, "/")
	return callLambdaFn("lint", lambdaContext{
		Registry:   sp[0],
		Username:   sp[1],
		Repository: sp[2],
		Linters:    linter,
	})
}
