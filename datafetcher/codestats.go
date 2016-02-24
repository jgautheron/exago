package datafetcher

import (
	"encoding/json"

	"github.com/exago/svc/repository"
)

var codeStatsCmd = &lambdaCmd{
	name:      "codestats",
	unMarshal: unMarshalCodeStats,
}

func GetCodeStats(repository string) (interface{}, error) {
	codeStatsCmd.ctxt = lambdaContext{
		Repository: repository,
	}
	return codeStatsCmd.Data()
}

func unMarshalCodeStats(l *lambdaCmd, b []byte) (interface{}, error) {
	var cs repository.CodeStats
	err := json.Unmarshal(b, &cs)
	return cs, err
}
