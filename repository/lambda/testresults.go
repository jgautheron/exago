package lambda

import (
	"encoding/json"

	"github.com/exago/svc/repository/model"
)

var testRunnerCmd = &cmd{
	name:      new(model.TestResults).Name(),
	unMarshal: unMarshalTestResults,
}

func GetTestResults(repository string) (model.RepositoryData, error) {
	testRunnerCmd.ctxt = context{
		Repository: repository,
	}
	return testRunnerCmd.Data()
}

func unMarshalTestResults(l *cmd, b []byte) (model.RepositoryData, error) {
	var tr model.TestResults
	err := json.Unmarshal(b, &tr)
	return tr, err
}
