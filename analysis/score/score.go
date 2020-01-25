package score

import (
	exago "github.com/jgautheron/exago/pkg"
	"github.com/sirupsen/logrus"
)

// Process triggers criterias evaluation, calling each evaluator in a goroutine
// We compute the weighted average based on the overall evaluator weights and scores
func Process(data exago.Data) (score float64, details []*exago.EvaluatorResponse) {
	eval := []CriteriaEvaluator{
		ThirdPartiesEvaluator(),
		CodeStatsEvaluator(),
		LintMessagesEvaluator(),
		TestCoverageEvaluator(),
		TestDurationEvaluator(),
		CheckListEvaluator(),
	}

	ch := make(chan *exago.EvaluatorResponse)
	for _, cr := range eval {
		go func(c CriteriaEvaluator) {
			c.Setup()
			ch <- c.Calculate(data)
		}(cr)
	}

	// Compute weighted average
	sw, avg := 0.0, 0.0
	res := []*exago.EvaluatorResponse{}
	for i := 0; i < len(eval); i++ {
		e := <-ch
		sw += e.Weight
		avg += e.Score * e.Weight

		logrus.WithField("score", e.Score).Debugf("[%s] score", e.Name)

		res = append(res, e)
	}

	if avg > 0 {
		avg /= sw
	}

	return avg, res
}
