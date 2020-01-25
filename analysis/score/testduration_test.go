package score_test

import (
	"testing"

	exago "github.com/jgautheron/exago/pkg"

	"github.com/jgautheron/exago/analysis/score"
)

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
		d := exago.Data{}
		d.ProjectRunner = getStubDuration(tt.duration)
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

func getStubDuration(duration []float64) exago.Results {
	tp := []exago.TestPackage{}
	for _, item := range duration {
		tp = append(tp, exago.TestPackage{ExecutionTime: item})
	}
	pr := exago.Results{}
	pr.CodeStats.Data = map[string]int{"loc": 123, "test": 123}
	pr.Test.Data = tp
	return pr
}
