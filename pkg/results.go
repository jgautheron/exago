package exago

import "simonwaldherr.de/go/golibs/xmath"

const (
	CodeStatsName    = "codestats"
	ThirdPartiesName = "thirdparties"
	TestCoverageName = "testcoverage"
	TestDurationName = "testduration"
)

type ChecklistItem struct {
	Category string `json:"category"`
	Desc     string `json:"desc"`
	Name     string `json:"name"`
}

type CoveragePackage struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Coverage float64 `json:"coverage"`
}

type TestPackage struct {
	Name          string     `json:"name"`
	ExecutionTime float64    `json:"executionTime"`
	Success       bool       `json:"success"`
	Tests         []TestFile `json:"tests"`
}

type TestFile struct {
	Name          string  `json:"name"`
	ExecutionTime float64 `json:"executionTime"`
	Passed        bool    `json:"passed"`
}

type Checklist struct {
	Failed []string `json:"failed"`
	Passed []string `json:"passed"`
}

// filename: []messages
type LinterResults map[string][]LinterFiles

type LinterFiles struct {
	Filename string          `json:"filename"`
	Messages []LinterMessage `json:"messages"`
}

type LinterMessage struct {
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Row      int    `json:"row"`
	Severity string `json:"severity"`
}

// Results received from the test runner.
type Results struct {
	Coverage struct {
		Label string `json:"label"`
		Data  struct {
			Packages []CoveragePackage `json:"packages"`
			Coverage float64           `json:"coverage"`
		} `json:"data"`
		RawOutput     string  `json:"rawOutput"`
		ExecutionTime float64 `json:"executionTime"`
	} `json:"coverage"`
	Download struct {
		Label         string      `json:"label"`
		Data          interface{} `json:"data"`
		RawOutput     string      `json:"rawOutput"`
		ExecutionTime float64     `json:"executionTime"`
	} `json:"download"`
	CodeStats struct {
		Label         string         `json:"label"`
		Data          map[string]int `json:"data"`
		RawOutput     string         `json:"rawOutput"`
		ExecutionTime float64        `json:"executionTime"`
	} `json:"codeStats"`
	Checklist struct {
		Label         string    `json:"label"`
		Data          Checklist `json:"data"`
		RawOutput     string    `json:"rawOutput"`
		ExecutionTime float64   `json:"executionTime"`
	} `json:"checklist"`
	Test struct {
		Label         string        `json:"label"`
		Data          []TestPackage `json:"data"`
		RawOutput     string        `json:"rawOutput"`
		ExecutionTime float64       `json:"executionTime"`
	} `json:"test"`
	ThirdParties struct {
		Label         string   `json:"label"`
		Data          []string `json:"data"`
		RawOutput     string   `json:"rawOutput"`
		ExecutionTime float64  `json:"executionTime"`
	} `json:"thirdParties"`
	Linters struct {
		Label         string        `json:"label"`
		Data          LinterResults `json:"data"`
		RawOutput     string        `json:"rawOutput"`
		ExecutionTime float64       `json:"executionTime"`
	} `json:"linters"`
}

// GetAvgTestDuration returns the average test duration.
func (t Results) GetMeanTestDuration() float64 {
	var duration []float64
	for _, pkg := range t.Test.Data {
		duration = append(duration, pkg.ExecutionTime)
	}
	if len(duration) == 0 {
		return 0
	}
	return xmath.Arithmetic(duration)
}

// GetAvgCodeCov returns the code coverage average.
func (t Results) GetMeanCodeCov() float64 {
	return t.Coverage.Data.Coverage
}
