package score

import (
	"fmt"
	"math"

	log "github.com/Sirupsen/logrus"
	"simonwaldherr.de/go/golibs/xmath"

	"github.com/exago/svc/repository/model"
)

const (
	drop      = 0.1
	dropSpeed = -0.1
)

type testDurationEvaluator struct {
	Evaluator
}

// TestDurationEvaluator measures a score based on test duration
func TestDurationEvaluator() CriteriaEvaluator {
	return &testDurationEvaluator{Evaluator{
		model.TestDurationName,
		"https://golang.org/pkg/testing/",
		"measures time taken for testing",
	}}
}

// Calculate overloads Evaluator/Calculate
func (te *testDurationEvaluator) Calculate(d model.Data) *model.EvaluatorResponse {
	t, cs := d.TestResults, d.CodeStats

	r := te.NewResponse(100, 3, "", nil)

	// If we have no tests, bypass the duration test
	if cs["Test"] == 0 {
		r.Score = 0
		r.Message = "no tests"

		return r
	}

	// Initialise values from test results
	var durations []float64
	for _, pkg := range t.Packages {
		durations = append(durations, pkg.ExecutionTime)
	}
	// Calculate mean values for execution time
	var duration float64
	if len(durations) > 0 {
		duration = xmath.Sum(durations)
	}

	log.WithFields(log.Fields{
		"duration (overall)": duration,
	}).Debugf("[%s] duration", model.TestDurationName)

	// Apply exponential growth formula
	r.Score = (100 + math.Exp(dropSpeed)) / (1 + math.Exp(dropSpeed-dropSpeed/drop)*duration)
	if r.Score > 100 {
		r.Score = 100
	}

	r.Message = fmt.Sprintf("tests took %.2fs", duration)

	return r
}
