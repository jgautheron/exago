package score

import "github.com/exago/svc/repository/model"

type codeStatsEvaluator struct {
	Evaluator
}

// CodeStatsEvaluator measures a score based on various metrics of code stats
// such as ratio LOC/CLOC and so on...
func CodeStatsEvaluator() CriteriaEvaluator {
	return &codeStatsEvaluator{Evaluator{
		model.CodeStatsName,
		"https://github.com/jgautheron/golocc",
		"counts lines of code, comments, functions, structs, imports etc in Go code",
	}}
}

// Calculate overloads Evaluator/Calculate
func (ce *codeStatsEvaluator) Calculate(d model.Data) *model.EvaluatorResponse {
	r := ce.NewResponse(100, 1, "", nil)
	cs := d.CodeStats
	ra := float64(cs["LOC"] / cs["NCLOC"])
	switch true {
	case ra > 1.4:
		r.Message = "more than 1.4 NCLOC"
	case ra > 1.3:
		r.Score = 75
		r.Message = "more than 1.3 NCLOC"
	case ra > 1.2:
		r.Score = 50
		r.Message = "more than 1.2 NCLOC"
	case ra > 1.1:
		r.Score = 25
		r.Message = "more than 1.1 NCLOC"
	case ra <= 1:
		r.Score = 0
		r.Message = "less or equal 1 NCLOC"
	}

	return r
}
