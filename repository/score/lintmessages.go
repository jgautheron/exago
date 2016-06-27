package score

import (
	"github.com/SimonWaldherr/golibs/xmath"
	log "github.com/Sirupsen/logrus"

	"github.com/exago/svc/repository/model"
)

func init() {
	l := make(map[string]*linter)
	l = map[string]*linter{
		"gofmt":       {0, 3, 0},
		"goimports":   {0, 2, 0},
		"golint":      {7, 1, 0},
		"dupl":        {5, 1.5, 0},
		"deadcode":    {0, 3, 0},
		"gocyclo":     {6, 2, 0},
		"vet":         {0, 2.5, 0},
		"vetshadow":   {0, 1, 0},
		"ineffassign": {0, 1, 0},
		"errcheck":    {6, 2, 0},
		"goconst":     {1.5, 1, 0},
		"gosimple":    {0, 1.5, 0},
		"staticcheck": {0, 1.5, 0},
	}

	Register(
		model.LintMessagesName,
		&LintMessagesEvaluator{
			Evaluator{100, 2, ""},
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
	weights := 0.0

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
			score := 100 - tmp
			weights += d.weight

			log.WithFields(log.Fields{
				"score":  score,
				"weight": d.weight,
			}).Debugf("[%s] score per linter", n)

			scores = append(scores, score*d.weight)
		}
	}

	if len(scores) > 0 {
		le.score = xmath.Sum(scores) / weights
	}

	log.WithFields(log.Fields{
		"score": le.score,
	}).Debugf("[%s] score", model.LintMessagesName)
}
