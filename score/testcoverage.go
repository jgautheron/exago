package score

import (
	"fmt"
	"math"

	log "github.com/Sirupsen/logrus"
	"simonwaldherr.de/go/golibs/xmath"

	"github.com/exago/svc/repository/model"
)

const coverageFactor = 0.1

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
func (te *testCoverageEvaluator) Calculate(p map[string]interface{}) *model.EvaluatorResponse {
	t := p[model.TestResultsName].(model.TestResults)
	cs := p[model.CodeStatsName].(model.CodeStats)

	r := te.NewResponse(100, 3, "", nil)

	// Initialise values from test results
	var cov []float64
	for _, pkg := range t.Packages {
		cov = append(cov, pkg.Coverage)
	}

	// Calculate mean for coverage
	var covMean float64
	if len(cov) > 0 {
		covMean = xmath.Geometric(cov)
	}

	log.WithFields(log.Fields{
		"coverage (geometric mean)": covMean,
	}).Debugf("[%s] coverage mean", model.TestCoverageName)

	// Apply exponential growth formula
	r.Score = covMean * math.Exp(coverageFactor)
	// Normalize to 100 if we go higher, this will probably never happen
	// but who knows...
	if r.Score > 100 {
		r.Score = 100
	}

	// Lines of code will impact the weight
	// We use a logarithm to calculate the factor on a base of 10
	r.Weight = math.Log10(float64(cs["LOC"]))

	switch true {
	case covMean > 0:
		r.Message = fmt.Sprintf("coverage is greater or equal to %.2f", covMean)
	case covMean == 0:
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
