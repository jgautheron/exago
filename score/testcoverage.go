package score

import (
	"fmt"
	"math"

	"github.com/hotolab/exago-svc/repository/model"
)

const coverageFactor = -0.095

type testCoverageEvaluator struct {
	Evaluator
}

// TestCoverageEvaluator measures a score based on test coverage
func TestCoverageEvaluator() CriteriaEvaluator {
	return &testCoverageEvaluator{Evaluator{
		model.TestCoverageName,
		"https://golang.org/pkg/testing/",
		"measures pourcentage of code covered by tests",
	}}
}

// Calculate overloads Evaluator/Calculate
func (te *testCoverageEvaluator) Calculate(d model.Data) *model.EvaluatorResponse {
	t, cs := d.ProjectRunner, d.CodeStats

	r := te.NewResponse(100, 3, "", nil)

	// Calculate mean for coverage
	covMean := t.GetMeanCodeCov()

	logger.WithField(
		"coverage (geometric mean)", covMean,
	).Debugf("[%s] coverage mean", model.TestCoverageName)

	// Apply logistic growth formula
	//
	// S = initial value = 1 (it starts from 1)
	// M = max value = 100 (the maximum score)
	// R = growth rate (negative)
	// V = value
	//
	// The equation is X = M/(S + [(M - S) * exp(R*V)])
	// this logistic model has two important parameters – a growth constant and a maximum size.
	r.Score = 100 / (1 + (100-1)*math.Exp(coverageFactor*covMean))

	// Lines of code will impact the weight
	// We use a logarithm to calculate the factor with a base10 of LOCs
	r.Weight = math.Log10(float64(cs["LOC"]))

	if covMean > 0 {
		r.Message = fmt.Sprintf("coverage is greater or equal to %.2f", covMean)
	} else {
		// If there are tests but we couldn't run them
		r.Score = 100
		r.Weight = 1
		r.Message = "some tests are available but could not be run, weight has been lowered"
		if cs["Test"] == 0 {
			r.Score = 0
			r.Weight = 3
			r.Message = "no tests"
		}
	}

	return r
}
