package score

import "github.com/exago/svc/repository/model"

func init() {
	Register(model.ImportsName, &ImportsEvaluator{Evaluator{100, 15, ""}})
}

// ImportsEvaluator measure a score based on various metrics of imports
// for now only the # of 3rd-party packages.
type ImportsEvaluator struct {
	Evaluator
}

// Calculate overloads Evaluator/Calculate
func (ie *ImportsEvaluator) Calculate(p map[string]interface{}) {
	imp := p[model.ImportsName].(model.Imports)
	tp := len(imp)
	switch true {
	case tp < 0:
	case tp < 4:
		ie.ScoreValue = 75
		ie.Msg = "less than 4"
	case tp < 6:
		ie.ScoreValue = 50
		ie.Msg = "less than 6"
	case tp < 8:
		ie.ScoreValue = 25
		ie.Msg = "less than 8"
	case tp > 8:
		ie.ScoreValue = 0
		ie.Msg = "more than 8"
	}
}
