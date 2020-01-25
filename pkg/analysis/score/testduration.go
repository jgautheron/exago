package score

import (
	"fmt"
	"math"

	"github.com/sirupsen/logrus"

	exago "github.com/jgautheron/exago/pkg"
)

const (
	// fastRate and slowRate are the two rate constants
	// expressed in reciprocal of the X (time) unit (inversed secs)
	fastRate = -0.1
	slowRate = -0.0008

	// percentFast is the fraction of the span (from initVal to plateau)
	// accounted for by the faster of the two components.
	percentFast = 28

	// initVal (Y0) is the Y value when X (time) is zero, represents the score
	initVal = 100

	// plateau is the Y value at infinite times, expressed in the same units as Y.
	// if duration is infinite score will be 0
	plateau = 0
)

type testDurationEvaluator struct {
	Evaluator
}

// TestDurationEvaluator measures a score based on test duration
func TestDurationEvaluator() CriteriaEvaluator {
	return &testDurationEvaluator{Evaluator{
		exago.TestDurationName,
		"https://golang.org/pkg/testing/",
		"measures time taken for testing",
	}}
}

// Calculate overloads Evaluator/Calculate
func (te *testDurationEvaluator) Calculate(d exago.Data) *exago.EvaluatorResponse {
	t, cs := d.Results, d.Results.CodeStats.Data

	r := te.NewResponse(100, 1.2, "", nil)

	// If we have no tests, bypass the duration test
	if cs["test"] == 0 {
		r.Score = 0
		r.Message = "no tests"

		return r
	}

	duration := t.GetMeanTestDuration()

	logrus.WithField(
		"duration (overall)", duration,
	).Debugf("[%s] duration", exago.TestDurationName)

	// A biphasic exponential decay or (two-phase) is used when the outcome is the result of
	// the sum of a fast and slow exponential decay.
	//
	// in this context test duration from 0 to 1s needs a different base rate
	// than longer duration. this is what we compute below.
	spanFast := initVal * percentFast * 0.01
	spanSlow := initVal * (initVal - percentFast) * 0.01

	r.Score = plateau + spanFast*math.Exp(fastRate*duration) + spanSlow*math.Exp(slowRate*duration)
	r.Message = fmt.Sprintf("tests took %.2fs", duration)

	return r
}
