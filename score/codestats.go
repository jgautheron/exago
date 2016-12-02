package score

import (
	"fmt"
	"math"

	"github.com/hotolab/exago-svc/repository/model"
)

const codeStatsFactor = -0.40

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
	r := ce.NewResponse(0, 1, "", nil)
	cs := d.ProjectRunner.CodeStats.Data
	ra := float64(cs["cloc"]) / float64(cs["loc"]) * 100

	r.Message = fmt.Sprintf("%d comments for %d lines of code", cs["cloc"], cs["loc"])

	if ra > 1 {
		r.Score = 100 / (1 + (30-1)*math.Exp(codeStatsFactor*ra))
	}

	return r
}
