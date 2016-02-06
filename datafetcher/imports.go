package datafetcher

import (
	"encoding/json"
	"strings"
)

func GetImports(repository string) (*json.RawMessage, error) {
	sp := strings.Split(repository, "/")
	return callLambdaFn("imports", lambdaContext{
		Registry:   sp[0],
		Username:   sp[1],
		Repository: sp[2],
	})
}
