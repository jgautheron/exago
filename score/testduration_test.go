package score_test

import (
	"testing"

	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/score"
)

func TestDurationNone(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = getStubDuration([]float64{0.2})
	d.CodeStats = map[string]int{"Test": 0}
	res := score.TestDurationEvaluator().Calculate(d)
	if res.Score != 0 {
		t.Error("The score should be 0")
	}
}

func TestDuration(t *testing.T) {
	var tests = []struct {
		duration []float64
		operator string
		expected float64
		desc     string
	}{
		{[]float64{0}, "=", 100, "The score should be 100"},
		{[]float64{5}, ">", 80, "5s is fast"},
		{[]float64{50}, "<", 70, "50s is still acceptable"},
		{[]float64{500}, "<", 50, "500s is long"},
		{[]float64{2000}, "<", 30, "2000s is unacceptable"},
	}

	for _, tt := range tests {
		d := model.Data{}
		d.ProjectRunner = getStubDuration(tt.duration)
		d.CodeStats = map[string]int{"Test": 123}
		res := score.TestDurationEvaluator().Calculate(d)

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

func getStubDuration(duration []float64) model.ProjectRunner {
	d := []model.Package{}
	for _, item := range duration {
		d = append(d, model.Package{ExecutionTime: item})
	}
	return model.ProjectRunner{
		Packages: d,
	}
}
