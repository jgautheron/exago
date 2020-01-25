package score

import (
	"fmt"
	"math"

	exago "github.com/jgautheron/exago/pkg"

	"github.com/sirupsen/logrus"
)

const coverageFactor = -0.11

type testCoverageEvaluator struct {
	Evaluator
}

// TestCoverageEvaluator measures a score based on test coverage
func TestCoverageEvaluator() CriteriaEvaluator {
	return &testCoverageEvaluator{Evaluator{
		exago.TestCoverageName,
		"https://golang.org/pkg/testing/",
		"measures pourcentage of code covered by tests",
	}}
}

// Calculate overloads Evaluator/Calculate
func (te *testCoverageEvaluator) Calculate(d exago.Data) *exago.EvaluatorResponse {
	t, cs := d.Results, d.Results.CodeStats.Data

	r := te.NewResponse(100, 3, "", nil)

	// Calculate mean for coverage
	covMean := t.GetMeanCodeCov()

	logrus.WithField(
		"coverage", covMean,
	).Debugf("[%s] coverage mean", exago.TestCoverageName)

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

	if covMean > 0 {
		r.Message = fmt.Sprintf("coverage is greater or equal to %.2f", covMean)
	} else {
		// If there are tests but we couldn't run them
		r.Score = 100
		r.Weight = 1
		r.Message = "some tests are available but could not be run, weight has been lowered"
		if cs["test"] == 0 {
			r.Score = 0
			r.Weight = 3
			r.Message = "no tests"
		}
	}

	return r
}
