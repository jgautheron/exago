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
	ExecutionTime float64    `json:"execution_time"`
	Success       bool       `json:"success"`
	Tests         []TestFile `json:"tests"`
}

type TestFile struct {
	Name          string  `json:"name"`
	ExecutionTime float64 `json:"execution_time"`
	Passed        bool    `json:"passed"`
}

type Checklist struct {
	Failed []ChecklistItem `json:"failed"`
	Passed []ChecklistItem `json:"passed"`
}

// Results received from the test runner.
type Results struct {
	Coverage struct {
		Label string `json:"label"`
		Data  struct {
			Packages []CoveragePackage `json:"packages"`
			Coverage float64           `json:"coverage"`
		} `json:"data"`
		RawOutput     string  `json:"raw_output"`
		ExecutionTime float64 `json:"execution_time"`
	} `json:"coverage"`
	Download struct {
		Label         string      `json:"label"`
		Data          interface{} `json:"data"`
		RawOutput     string      `json:"raw_output"`
		ExecutionTime float64     `json:"execution_time"`
	} `json:"download"`
	CodeStats struct {
		Label         string         `json:"label"`
		Data          map[string]int `json:"data"`
		RawOutput     string         `json:"raw_output"`
		ExecutionTime float64        `json:"execution_time"`
	} `json:"codeStats"`
	Checklist struct {
		Label         string    `json:"label"`
		Data          Checklist `json:"data"`
		RawOutput     string    `json:"raw_output"`
		ExecutionTime float64   `json:"execution_time"`
	} `json:"checklist"`
	Test struct {
		Label         string        `json:"label"`
		Data          []TestPackage `json:"data"`
		RawOutput     string        `json:"raw_output"`
		ExecutionTime float64       `json:"execution_time"`
	} `json:"test"`
	ThirdParties struct {
		Label         string   `json:"label"`
		Data          []string `json:"data"`
		RawOutput     string   `json:"raw_output"`
		ExecutionTime float64  `json:"execution_time"`
	} `json:"thirdParties"`
	Linters struct {
		Label         string                                         `json:"label"`
		Data          map[string]map[string][]map[string]interface{} `json:"data"`
		RawOutput     string                                         `json:"raw_output"`
		ExecutionTime float64                                        `json:"execution_time"`
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
