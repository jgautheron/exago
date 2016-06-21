package score

import (
	"github.com/SimonWaldherr/golibs/xmath"
	"github.com/exago/svc/repository/model"
)

func init() {
	Register(model.LintMessagesName, &LintMessagesEvaluator{Evaluator{100, .25, ""}})
}

// LintMessagesEvaluator measure a score based on the output of gometalinter
type LintMessagesEvaluator struct {
	Evaluator
}

// Calculate overloads Evaluator/Calculate
func (ce *LintMessagesEvaluator) Calculate(p map[string]interface{}) {
	var linters = map[string]*struct{ threshold, warnings float64 }{
		"gocyclo":  {.05, 0},
		"golint":   {.05, 0},
		"errcheck": {.05, 0},
	}

	lm := p[model.LintMessagesName].(model.LintMessages)
	cs := p[model.CodeStatsName].(model.CodeStats)

	// Loop over messages
	for _, m := range lm {
		for ln, lr := range m {
			if l, ok := linters[ln]; ok {
				for _, a := range lr {
					if a["severity"].(string) == "warning" {
						l.warnings++
					}
				}
			}
		}
	}

	scores := []float64{}
	// Computing score
	for _, d := range linters {
		tmp := 100 * d.warnings / float64(cs["LOC"])
		if tmp > d.threshold {
			score := 100 - tmp
			scores = append(scores, score)
		}
	}

	ce.score = xmath.Arithmetic(scores)
}
