package repository

// TestResults received from the test runner.
type TestResults struct {
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

type Imports []string
type CodeStats map[string]int
type LintMessages map[string]map[string][]map[string]interface{}
