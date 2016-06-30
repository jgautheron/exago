package score

import (
	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/repository/model"
)

// Process triggers criterias evaluation
func Process(params map[string]interface{}) (score float64, details []*model.EvaluatorResponse) {
	eval := []CriteriaEvaluator{
		ImportsEvaluator(),
		CodeStatsEvaluator(),
		LintMessagesEvaluator(),
		TestCoverageEvaluator(),
		TestDurationEvaluator(),
		CheckListEvaluator(),
	}

	ch := make(chan *model.EvaluatorResponse)
	for _, cr := range eval {
		go func(c CriteriaEvaluator) {
			c.Setup()
			ch <- c.Calculate(params)
		}(cr)
	}

	// Compute weighted average
	sw, avg := 0.0, 0.0
	res := []*model.EvaluatorResponse{}
	for i := 0; i < len(eval); i++ {
		e := <-ch
		sw += e.Weight
		avg += e.Score * e.Weight

		log.WithFields(log.Fields{
			"score": e.Score,
		}).Debugf("[%s] score", e.Name)

		res = append(res, e)
	}

	if avg > 0 {
		avg /= sw
	}

	return avg, res
}
