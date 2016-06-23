package score

import "github.com/exago/svc/repository/model"

func init() {
	Register(model.CodeStatsName, &CodeStatsEvaluator{Evaluator{100, 10, ""}})
}

// CodeStatsEvaluator measure a score based on various metrics of code stats
// such as ratio LOC/CLOC and so on...
type CodeStatsEvaluator struct {
	Evaluator
}

// Calculate overloads Evaluator/Calculate
func (ce *CodeStatsEvaluator) Calculate(p map[string]interface{}) {
	cs := p[model.CodeStatsName].(model.CodeStats)
	ra := float64(cs["LOC"] / cs["NCLOC"])
	switch true {
	case ra > 1.4:
		ce.Msg = "more than 1.4"
	case ra > 1.3:
		ce.ScoreValue = 75
		ce.Msg = "more than 1.3"
	case ra > 1.2:
		ce.ScoreValue = 50
		ce.Msg = "more than 1.2"
	case ra > 1.1:
		ce.ScoreValue = 25
		ce.Msg = "more than 1.1"
	case ra <= 1:
		ce.ScoreValue = 0
		ce.Msg = "less or equal 1"
	}
}
