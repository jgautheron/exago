package datafetcher

import (
	"encoding/json"

	"github.com/exago/svc/repository"
)

var lintCmd = &lambdaCmd{
	name:      "lint",
	unMarshal: unMarshalLint,
}

func GetLintMessages(repository, linter string) (interface{}, error) {
	lintCmd.ctxt = lambdaContext{
		Repository: repository,
		Linter:     linter,
	}
	return lintCmd.Data()
}

func unMarshalLint(l *lambdaCmd, b []byte) (interface{}, error) {
	var lnt repository.LintMessages
	err := json.Unmarshal(b, &lnt)
	return lnt, err
}
