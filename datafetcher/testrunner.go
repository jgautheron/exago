package datafetcher

import "encoding/json"

type testResults struct {
	Checklist struct {
		Failed []struct {
			Category string `json:"Category"`
			Desc     string `json:"Desc"`
			Name     string `json:"Name"`
		} `json:"Failed"`
		Passed []struct {
			Category string `json:"Category"`
			Desc     string `json:"Desc"`
			Name     string `json:"Name"`
		} `json:"Passed"`
	} `json:"checklist"`
	Packages []struct {
		Coverage      float64 `json:"coverage"`
		ExecutionTime float64 `json:"execution_time"`
		Name          string  `json:"name"`
		Success       bool    `json:"success"`
	} `json:"packages"`
}

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
	var tr testResults
	err := json.Unmarshal(b, &tr)
	return tr, err
}
