package datafetcher

import "strings"

func GetCodeStats(repository string) (LambdaResponse, error) {
	sp := strings.Split(repository, "/")
	return callLambdaFn("loc", lambdaContext{
		Registry:   sp[0],
		Username:   sp[1],
		Repository: sp[2],
	})
}
