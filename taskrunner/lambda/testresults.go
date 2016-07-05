package lambda

import (
	"encoding/json"

	"github.com/exago/svc/repository/model"
)

var testRunnerCmd = &cmd{
	name:      model.TestResultsName,
	unMarshal: unMarshalTestResults,
}

func (l Runner) FetchTestResults() (interface{}, error) {
	testRunnerCmd.ctxt = context{
		Repository: l.Repository,
		Cleanup:    l.ShouldCleanup,
	}
	return testRunnerCmd.Data()
}

func unMarshalTestResults(l *cmd, b []byte) (interface{}, error) {
	var tr model.TestResults
	err := json.Unmarshal(b, &tr)
	return tr, err
}
