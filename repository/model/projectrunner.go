package model

import "simonwaldherr.de/go/golibs/xmath"

const (
	ProjectRunnerName = "projectrunner"
	ThirdPartiesName  = "thirdparties"
	TestCoverageName  = "testcoverage"
	TestDurationName  = "testduration"
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

// ProjectRunner received from the test runner.
type ProjectRunner struct {
	Coverage struct {
		Label string `json:"label"`
		Data  struct {
			Packages []CoveragePackage `json:"packages"`
			Coverage float64           `json:"coverage"`
		} `json:"data"`
		RawOutput     string      `json:"raw_output"`
		ExecutionTime float64     `json:"execution_time"`
		Error         interface{} `json:"error"`
	} `json:"coverage"`
	Download struct {
		Label         string      `json:"label"`
		Data          interface{} `json:"data"`
		RawOutput     string      `json:"raw_output"`
		ExecutionTime float64     `json:"execution_time"`
		Error         interface{} `json:"error"`
	} `json:"download"`
	Goprove struct {
		Label         string      `json:"label"`
		Data          Checklist   `json:"data"`
		RawOutput     string      `json:"raw_output"`
		ExecutionTime float64     `json:"execution_time"`
		Error         interface{} `json:"error"`
	} `json:"goprove"`
	Test struct {
		Label         string        `json:"label"`
		Data          []TestPackage `json:"data"`
		RawOutput     string        `json:"raw_output"`
		ExecutionTime float64       `json:"execution_time"`
		Error         interface{}   `json:"error"`
	} `json:"test"`
	Thirdparties struct {
		Label         string      `json:"label"`
		Data          []string    `json:"data"`
		RawOutput     string      `json:"raw_output"`
		ExecutionTime float64     `json:"execution_time"`
		Error         interface{} `json:"error"`
	} `json:"thirdparties"`
}

// GetAvgTestDuration returns the average test duration.
func (t ProjectRunner) GetMeanTestDuration() float64 {
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
func (t ProjectRunner) GetMeanCodeCov() float64 {
	return t.Coverage.Data.Coverage
}
