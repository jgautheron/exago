package lambda

import (
	"encoding/json"
	"strings"

	"github.com/hotolab/exago-svc/repository/model"
)

var (
	// DefaultLinters ran by default in Lambda.
	DefaultLinters = []string{
		"deadcode", "dupl", "gas", "goconst", "gocyclo", "gofmt",
		"golint", "gosimple", "ineffassign", "staticcheck", "vet", "vetshadow",
	}
)

func (l Runner) FetchLintMessages() (model.LintMessages, error) {
	lintCmd := &cmd{
		name:      model.LintMessagesName,
		unMarshal: unMarshalLint,
	}
	lintCmd.ctxt = context{
		Repository: l.Repository,
		Cleanup:    l.ShouldCleanup,
		Linters:    strings.Join(DefaultLinters, ","),
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
