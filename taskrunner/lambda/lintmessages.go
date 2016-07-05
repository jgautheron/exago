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

func (l Runner) FetchLintMessages(linters []string) (model.LintMessages, error) {
	lintCmd.ctxt = context{
		Repository: l.Repository,
		Cleanup:    l.ShouldCleanup,
		Linters:    strings.Join(linters, ","),
	}
	d, err := lintCmd.Data()
	if err != nil {
		return nil, err
	}
	return d.(model.LintMessages), nil
}

func unMarshalLint(l *cmd, b []byte) (interface{}, error) {
	var lnt model.LintMessages
	err := json.Unmarshal(b, &lnt)
	return lnt, err
}
