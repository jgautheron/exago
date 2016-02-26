package datafetcher

import (
	"encoding/json"

	"github.com/exago/svc/repository"
)

var testRunnerCmd = &lambdaCmd{
	name:      "testrunner",
	unMarshal: unMarshalTestResults,
}

func GetTestResults(repository string) (interface{}, error) {
	testRunnerCmd.ctxt = lambdaContext{
		Repository: repository,
	}
	return testRunnerCmd.Data()
}

func unMarshalTestResults(l *lambdaCmd, b []byte) (interface{}, error) {
	var tr repository.TestResults
	err := json.Unmarshal(b, &tr)
	return tr, err
}
