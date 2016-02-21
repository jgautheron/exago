package datafetcher

import "encoding/json"

type codeStats map[string]int

var codeStatsCmd = &lambdaCmd{
	name:      "loc",
	unMarshal: unMarshalCodeStats,
}

func GetCodeStats(repository string) (interface{}, error) {
	codeStatsCmd.ctxt = lambdaContext{
		Repository: repository,
	}
	return codeStatsCmd.Data()
}

func unMarshalCodeStats(l *lambdaCmd, b []byte) (interface{}, error) {
	var cs codeStats
	err := json.Unmarshal(b, &cs)
	return cs, err
}
