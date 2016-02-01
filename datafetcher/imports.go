package datafetcher

import "strings"

func GetImports(repository string) (LambdaResponse, error) {
	sp := strings.Split(repository, "/")
	return callLambdaFn("imports", lambdaContext{
		Registry:   sp[0],
		Username:   sp[1],
		Repository: sp[2],
	})
}
