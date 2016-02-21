package datafetcher

import "encoding/json"

type lint map[string]map[string][]map[string]interface{}

var lintCmd = &lambdaCmd{
	name:      "lint",
	unMarshal: unMarshalLint,
}

func GetLint(repository, linter string) (interface{}, error) {
	lintCmd.ctxt = lambdaContext{
		Repository: repository,
		Linter:     linter,
	}
	return lintCmd.Data()
}

func unMarshalLint(l *lambdaCmd, b []byte) (interface{}, error) {
	var lnt lint
	err := json.Unmarshal(b, &lnt)
	return lnt, err
}
