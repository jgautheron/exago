package lambda

import (
	"encoding/json"
	"strings"

	"github.com/exago/svc/repository/model"
)

var lintCmd = &cmd{
	name:      new(model.LintMessages).Name(),
	unMarshal: unMarshalLint,
}

func GetLintMessages(repository string, linters []string) (model.RepositoryData, error) {
	lintCmd.ctxt = context{
		Repository: repository,
		Linters:    strings.Join(linters, ","),
	}
	return lintCmd.Data()
}

func unMarshalLint(l *cmd, b []byte) (model.RepositoryData, error) {
	var lnt model.LintMessages
	err := json.Unmarshal(b, &lnt)
	return lnt, err
}
