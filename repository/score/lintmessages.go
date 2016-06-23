package score

import (
	"fmt"

	"github.com/SimonWaldherr/golibs/xmath"
	log "github.com/Sirupsen/logrus"

	"github.com/exago/svc/repository/model"
)

func init() {
	l := make(map[string]*linter)
	l = map[string]*linter{
		"gofmt":       {0, 35, 0},
		"goimports":   {0, 20, 0},
		"golint":      {7, 10, 0},
		"dupl":        {5, 15, 0},
		"deadcode":    {0, 30, 0},
		"gocyclo":     {6, 20, 0},
		"vet":         {0, 25, 0},
		"vetshadow":   {0, 10, 0},
		"ineffassign": {0, 10, 0},
		"errcheck":    {6, 20, 0},
		"goconst":     {1.5, 10, 0},
		"gosimple":    {0, 15, 0},
		"staticcheck": {0, 15, 0},
	}

	Register(
		model.LintMessagesName,
		&LintMessagesEvaluator{
			Evaluator{100, 25, ""},
			l,
		},
	)
}

type linter struct {
	threshold float64
	weight    float64
	warnings  float64
}

// LintMessagesEvaluator measure a score based on the output of gometalinter
type LintMessagesEvaluator struct {
	Evaluator
	linters map[string]*linter
}

// Calculate overloads Evaluator/Calculate
func (le *LintMessagesEvaluator) Calculate(p map[string]interface{}) {
	le.normalizeWeightSerie()

	lm := p[model.LintMessagesName].(model.LintMessages)
	cs := p[model.CodeStatsName].(model.CodeStats)

	// Loop over messages
	for _, m := range lm {
		for ln, lr := range m {
			if l, ok := le.linters[ln]; ok {
				for _, a := range lr {
					// Count the number of warnings
					if a["severity"].(string) == "warning" {
						l.warnings++
					}
				}
			}
		}
	}

	// Compute score
	scores := []float64{}
	for n, d := range le.linters {
		tmp := 100 * d.warnings / float64(cs["LOC"])
		log.WithFields(log.Fields{
			"defect ratio": tmp,
			"threshold":    d.threshold,
			"warnings":     d.warnings,
			"loc":          cs["LOC"],
			"weight":       d.weight,
		}).Debugf("[%s] threshold vs ratio", n)
		if tmp > d.threshold {
			// Compute weighted score
			score := (100 - tmp) * d.weight
			le.Msg += fmt.Sprintf("[%s] = %.2f", n, score/d.weight)
			scores = append(scores, score)
		}
	}

	if len(scores) > 0 {
		le.ScoreValue -= xmath.Sum(scores) / 100
	}
}

func (le *LintMessagesEvaluator) normalizeWeightSerie() {
	w := []float64{}
	for _, l := range le.linters {
		w = append(w, l.weight)
	}
	// Make the sum of all weights, and calculate the scaling factor
	factor := xmath.Sum(w) / 100
	// Loop against each linter and fix the new weight
	for n, l := range le.linters {
		if l.weight > 0 {
			l.weight /= factor
			log.WithField("weight", l.weight).Debugf("[%s] Calculated weight", n)
		}
	}
}
