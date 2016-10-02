package score

import (
	"fmt"
	"math"

	"github.com/hotolab/exago-svc/repository/model"
)

type thirdPartiesEvaluator struct {
	Evaluator
}

// ThirdPartiesEvaluator measures a score based on various metrics of imports
// for now only the # of 3rd-party packages.
func ThirdPartiesEvaluator() CriteriaEvaluator {
	return &thirdPartiesEvaluator{Evaluator{
		model.ThirdPartiesName,
		"https://github.com/jgautheron/gogetimports",
		"counts the number of third party libraries",
	}}
}

// Calculate overloads Evaluator/Calculate
func (ie *thirdPartiesEvaluator) Calculate(d model.Data) *model.EvaluatorResponse {
	// Declare rates here, since Go cannot accept maps as constants :/
	rates := map[int]float64{
		1: -1,
		2: -0.8,
		3: -0.6,
		4: -0.35,
		5: -0.25,
		6: -0.18,
	}

	imp, cs := d.ProjectRunner.Thirdparties.Data, d.CodeStats

	imps := float64(len(imp))
	r := ie.NewResponse(100, 1.5, "", nil)

	loc := float64(cs["LOC"])

	// We simply compute the power of 10 using log10 and floor
	// and retrieve the rate by associating the power to an index position
	l10 := math.Floor(math.Log10(loc))
	rate := rates[5]

	// If we can't find the rate, fallback to the lowest rate
	// This will unlikely happen over 1,000,000 LOC
	if val, ok := rates[int(l10)]; ok {
		rate = val
	}

	// Compute the exponential decay
	r.Score = 100 * math.Exp(rate*(imps/math.Log(loc)))
	r.Message = fmt.Sprintf("%d third-party package(s)", int(imps))

	return r
}
