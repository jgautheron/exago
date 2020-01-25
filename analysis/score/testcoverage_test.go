package score_test

import (
	"testing"

	exago "github.com/jgautheron/exago/pkg"

	"simonwaldherr.de/go/golibs/xmath"

	"github.com/jgautheron/exago/analysis/score"
)

func TestCoverageNone(t *testing.T) {
	d := exago.Data{}
	d.ProjectRunner = getStubCoverage([]float64{})
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
		d := exago.Data{}
		d.ProjectRunner = getStubCoverage(tt.coverage)
		res := score.TestCoverageEvaluator().Calculate(d)

		switch tt.operator {
		case "<":
			if res.Score > tt.expected {
				t.Errorf("Wrong score %s: %d is not > to %d", tt.desc, res.Score, tt.expected)
			}
		case ">":
			if res.Score < tt.expected {
				t.Errorf("Wrong score %s: %d is not < to %d", tt.desc, res.Score, tt.expected)
			}
		case "=":
			if res.Score != tt.expected {
				t.Errorf("Wrong score %s: %d is not = to %d", tt.desc, res.Score, tt.expected)
			}
		}
	}
}

func getStubCoverage(coverage []float64) exago.Results {
	pr := exago.Results{}
	pr.Coverage.Data.Coverage = xmath.Arithmetic(coverage)
	pr.CodeStats.Data = map[string]int{"loc": 123}
	return pr
}
