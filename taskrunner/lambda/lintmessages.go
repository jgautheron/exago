package lambda

import (
	"encoding/json"
	"strings"

	"github.com/exago/svc/repository/model"
)

var lintCmd = &cmd{
	name:      model.LintMessagesName,
	unMarshal: unMarshalLint,
}

func (l Runner) FetchLintMessages(linters []string) (interface{}, error) {
	lintCmd.ctxt = context{
		Repository: l.repository,
		Linters:    strings.Join(linters, ","),
		Cleanup:    l.shouldCleanup,
	}
	return lintCmd.Data()
}

func unMarshalLint(l *cmd, b []byte) (interface{}, error) {
	var lnt model.LintMessages
	err := json.Unmarshal(b, &lnt)
	return lnt, err
}
