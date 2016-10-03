package model

import "testing"

func TestGotMeanDuration(t *testing.T) {
	tr := getMockProjectRunner()
	if tr.GetMeanTestDuration() != 0.05 {
		t.Errorf("Got the wrong mean test duration")
	}
}

func TestGotNullMeanDuration(t *testing.T) {
	tr := ProjectRunner{}
	if tr.GetMeanTestDuration() != 0 {
		t.Errorf("Got the wrong mean test duration")
	}
}

func TestGotMeanCoverage(t *testing.T) {
	tr := getMockProjectRunner()
	if tr.GetMeanCodeCov() != 20.799 {
		t.Errorf("Got the wrong mean test coverage")
	}
}

func TestGotNullMeanCoverage(t *testing.T) {
	tr := ProjectRunner{}
	if tr.GetMeanCodeCov() != 0 {
		t.Errorf("Got the wrong mean test coverage")
	}
}

func getMockProjectRunner() ProjectRunner {
	tr := ProjectRunner{}
	tr.Coverage.Data.Coverage = 20.799
	tr.Test.Data = append(tr.Test.Data, struct {
		Name          string  `json:"name"`
		ExecutionTime float64 `json:"execution_time"`
		Success       bool    `json:"success"`
		Tests         []struct {
			Name          string  `json:"name"`
			ExecutionTime float64 `json:"execution_time"`
			Passed        bool    `json:"passed"`
		} `json:"tests"`
	}{
		ExecutionTime: 0.05,
	})
	return tr
}
