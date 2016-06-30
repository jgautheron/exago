package score

import (
	log "github.com/Sirupsen/logrus"

	"github.com/exago/svc/repository/model"
)

type testResultsEvaluator struct {
	Evaluator
}

// TestResultsEvaluator measures a score based on test results
func TestResultsEvaluator() CriteriaEvaluator {
	return &testResultsEvaluator{Evaluator{
		model.TestResultsName,
		"https://golang.org/pkg/testing/",
		"automated testing of Go packages",
	}}
}

// Calculate overloads Evaluator/Calculate
func (te *testResultsEvaluator) Calculate(p map[string]interface{}) *model.EvaluatorResponse {
	t := p[model.TestResultsName].(model.TestResults)
	cs := p[model.CodeStatsName].(model.CodeStats)

	r := te.NewResponse(100, 3, "", nil)

	// Initialise values from test results
	var cov, duration []float64
	for _, pkg := range t.Packages {
		cov = append(cov, pkg.Coverage)
		duration = append(duration, pkg.ExecutionTime)
	}

	// Calculate mean values for both code coverage and execution time
	var covMean, durationMean float64 = 0, 0
	if len(cov) > 0 {
		for _, v := range cov {
			covMean += v
		}
		covMean /= float64(len(cov))
	}
	if len(duration) > 0 {
		for _, v := range duration {
			durationMean += v
		}
		durationMean /= float64(len(duration))
	}

	log.WithFields(log.Fields{
		"coverage (mean)": covMean,
		"duration (mean)": durationMean,
	}).Debugf("[%s] coverage and duration mean", model.TestResultsName)

	switch true {
	case covMean > 70:
		r.Message = "coverage is greather than 70"
	case covMean > 60:
		r.Score = 80
		r.Message = "coverage is greater than 60"
	case covMean > 40:
		r.Score = 60
		r.Message = "coverage is greater than 40"
	case covMean == 0:
		// If there are tests but we couldn't run them, do not  points
		r.Score = 100
		r.Weight = 1
		r.Message = "some tests are available but could not be run, weight has been lowered"
		if cs["Test"] == 0 {
			r.Score = 0
			r.Weight = 3
			r.Message = "no tests"
		}
	}

	// Fast test suites are important
	switch true {
	case durationMean == 0:
		// Do nothing
	case durationMean > 10:
		r.Weight = 2.5
	case durationMean < 2:
		r.Weight = 2
	}

	return r
}
