package lambda

import (
	"encoding/json"

	"github.com/hotolab/exago-svc/repository/model"
)

var testRunnerCmd = &cmd{
	name:      model.TestResultsName,
	unMarshal: unMarshalTestResults,
}

func (l Runner) FetchTestResults() (tr model.TestResults, err error) {
	testRunnerCmd.ctxt = context{
		Repository: l.Repository,
		Cleanup:    l.ShouldCleanup,
	}
	d, err := testRunnerCmd.Data()
	if err != nil {
		return tr, err
	}
	return d.(model.TestResults), nil
}

func unMarshalTestResults(l *cmd, b []byte) (interface{}, error) {
	var tr model.TestResults
	err := json.Unmarshal(b, &tr)
	return tr, err
}
