package score

import (
	"fmt"
	"math"

	"github.com/exago/svc/repository/model"
)

const codeStatsFactor = -0.13

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
func (ce *codeStatsEvaluator) Calculate(p map[string]interface{}) *model.EvaluatorResponse {
	r := ce.NewResponse(100, 1, "", nil)
	cs := p[model.CodeStatsName].(model.CodeStats)
	ra := float64(cs["CLOC"]) / float64(cs["LOC"])
	r.Message = fmt.Sprintf("%d comments for %d lines of code", cs["CLOC"], cs["LOC"])

	r.Score = 100 / (1 + (100-1)*math.Exp(codeStatsFactor*(ra*100)))

	return r
}
