package model

import (
	"testing"

	"simonwaldherr.de/go/golibs/xmath"
)

func TestGotMeanDuration(t *testing.T) {
	tr := getMockProjectRunner()
	if tr.GetAvgTestDuration() != xmath.Arithmetic([]float64{0.01, 0.05}) {
		t.Errorf("Got the wrong mean test duration")
	}
}

func TestGotNullMeanDuration(t *testing.T) {
	tr := ProjectRunner{}
	if tr.GetAvgTestDuration() != 0 {
		t.Errorf("Got the wrong mean test duration")
	}
}

func TestGotMeanCoverage(t *testing.T) {
	tr := getMockProjectRunner()
	if tr.GetAvgCodeCov() != xmath.Geometric([]float64{20.799, 80.0001}) {
		t.Errorf("Got the wrong mean test coverage")
	}
}

func TestGotNullMeanCoverage(t *testing.T) {
	tr := ProjectRunner{}
	if tr.GetAvgCodeCov() != 0 {
		t.Errorf("Got the wrong mean test coverage")
	}
}

func getMockProjectRunner() ProjectRunner {
	tr := ProjectRunner{}
	tr.Packages = append(tr.Packages, struct {
		Coverage      float64 `json:"coverage"`
		ExecutionTime float64 `json:"execution_time"`
		Name          string  `json:"name"`
		Success       bool    `json:"success"`
		Tests         []struct {
			Name          string  `json:"name"`
			ExecutionTime float64 `json:"execution_time"`
			Passed        bool    `json:"passed"`
		} `json:"tests"`
	}{
		Coverage:      20.799,
		ExecutionTime: 0.01,
	})
	tr.Packages = append(tr.Packages, struct {
		Coverage      float64 `json:"coverage"`
		ExecutionTime float64 `json:"execution_time"`
		Name          string  `json:"name"`
		Success       bool    `json:"success"`
		Tests         []struct {
			Name          string  `json:"name"`
			ExecutionTime float64 `json:"execution_time"`
			Passed        bool    `json:"passed"`
		} `json:"tests"`
	}{
		Coverage:      80.0001,
		ExecutionTime: 0.05,
	})
	return tr
}
