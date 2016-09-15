package score_test

import (
	"testing"

	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/score"
)

func TestCoverageNone(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = getStubCoverage([]float64{})
	d.CodeStats = map[string]int{"Test": 0}
	res := score.TestCoverageEvaluator().Calculate(d)
	if res.Score != 0 {
		t.Error("The score should be 0")
	}
}

func TestCoverage(t *testing.T) {
	var tests = []struct {
		coverage []float64
		operator string
		expected float64
		desc     string
	}{
		{[]float64{100}, ">", 99, "The score should be around 100"},
		{[]float64{80, 50, 80}, ">", 80, "Pretty good coverage"},
		{[]float64{10, 20, 30}, "<", 20, "Bad!"},
		{[]float64{50}, ">", 50, "50% is not bad"},
	}

	for _, tt := range tests {
		d := model.Data{}
		d.ProjectRunner = getStubCoverage(tt.coverage)
		d.CodeStats = map[string]int{"Test": 123}
		res := score.TestCoverageEvaluator().Calculate(d)

		switch tt.operator {
		case "<":
			if res.Score > tt.expected {
				t.Errorf("Wrong score %s", tt.desc)
			}
		case ">":
			if res.Score < tt.expected {
				t.Errorf("Wrong score %s", tt.desc)
			}
		case "=":
			if res.Score != tt.expected {
				t.Errorf("Wrong score %s", tt.desc)
			}
		}
	}
}

func getStubCoverage(coverage []float64) model.ProjectRunner {
	d := []model.Package{}
	for _, item := range coverage {
		d = append(d, model.Package{Coverage: item})
	}
	return model.ProjectRunner{
		Packages: d,
	}
}
