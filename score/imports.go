package score

import "github.com/exago/svc/repository/model"

type importsEvaluator struct {
	Evaluator
}

// ImportsEvaluator measures a score based on various metrics of imports
// for now only the # of 3rd-party packages.
func ImportsEvaluator() CriteriaEvaluator {
	return &importsEvaluator{Evaluator{
		model.ImportsName,
		"https://github.com/jgautheron/gogetimports",
		"counts the number of third party libraries",
	}}
}

// Calculate overloads Evaluator/Calculate
func (ie *importsEvaluator) Calculate(d model.Data) *model.EvaluatorResponse {
	imp := d.Imports
	tp := len(imp)
	r := ie.NewResponse(100, 1.5, "", nil)

	switch true {
	case tp < 0:
	case tp < 4:
		r.Score = 75
		r.Message = "less than 4 third-party package(s)"
	case tp < 6:
		r.Score = 50
		r.Message = "less than 6 third-party package(s)"
	case tp < 8:
		r.Score = 25
		r.Message = "less than 8 third-party package(s)"
	case tp > 8:
		r.Score = 0
		r.Message = "more than 8 third-party package(s)"
	}

	return r
}
