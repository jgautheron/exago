package model

const (
	ImportsName       = "imports"
	CodeStatsName     = "codestats"
	LintMessagesName  = "lintmessages"
	TestResultsName   = "testresults"
	ScoreName         = "score"
	ExecutionTimeName = "execution_time"
	LastUpdateName    = "last_update"
)

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
		Tests         []struct {
			Name          string  `json:"name"`
			ExecutionTime float64 `json:"execution_time"`
			Passed        bool    `json:"passed"`
		} `json:"tests"`
	} `json:"packages"`
	ExecutionTime struct {
		Goget   string `json:"goget,omitempty"`
		Goprove string `json:"goprove"`
		Gotest  string `json:"gotest"`
	} `json:"execution_time"`
	RawOutput struct {
		Goget  string `json:"goget"`
		Gotest string `json:"gotest"`
	} `json:"raw_output"`
	Errors struct {
		Goget  string `json:"goget"`
		Gotest string `json:"gotest"`
	} `json:"errors"`
}

type Imports []string
type CodeStats map[string]int
type LintMessages map[string]map[string][]map[string]interface{}

type Score struct {
	Value   float64  `json:"value"`
	Rank    string   `json:"rank"`
	Details []string `json:"details"`
}

func (t TestResults) GetAvgTestDuration() float64 {
	var duration []float64
	for _, pkg := range t.Packages {
		duration = append(duration, pkg.ExecutionTime)
	}

	var durationMean float64 = 0
	if len(duration) > 0 {
		for _, v := range duration {
			durationMean += v
		}
		durationMean /= float64(len(duration))
	}

	return durationMean
}

func (t TestResults) GetAvgCodeCov() float64 {
	var cov []float64
	for _, pkg := range t.Packages {
		cov = append(cov, pkg.Coverage)
	}

	var covMean float64 = 0
	if len(cov) > 0 {
		for _, v := range cov {
			covMean += v
		}
		covMean /= float64(len(cov))
	}

	return covMean
}
