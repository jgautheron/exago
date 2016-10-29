package lambda

import (
	"encoding/json"

	"github.com/hotolab/exago-svc/repository/model"
)

func (l Runner) FetchProjectRunner() (tr model.ProjectRunner, err error) {
	runnerCmd := &cmd{
		name:      model.ProjectRunnerName,
		unMarshal: unMarshalProjectRunner,
	}
	runnerCmd.ctxt = context{
		Repository: l.Repository,
		Cleanup:    l.ShouldCleanup,
	}
	d, err := runnerCmd.Data()
	if err != nil {
		return tr, err
	}
	return d.(model.ProjectRunner), nil
}

func unMarshalProjectRunner(l *cmd, b []byte) (interface{}, error) {
	var tr model.ProjectRunner
	err := json.Unmarshal(b, &tr)
	return tr, err
}
