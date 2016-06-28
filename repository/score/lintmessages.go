package score

import (
	"github.com/SimonWaldherr/golibs/xmath"
	log "github.com/Sirupsen/logrus"

	"github.com/exago/svc/repository/model"
)

type linter struct {
	threshold float64
	weight    float64
	warnings  float64
	url       string
	desc      string
}

type lintMessagesEvaluator struct {
	Evaluator
	linters map[string]*linter
}

// LintMessagesEvaluator measures a score based on the output of gometalinter
func LintMessagesEvaluator() CriteriaEvaluator {
	return &lintMessagesEvaluator{Evaluator{
		model.LintMessagesName,
		"https://github.com/alecthomas/gometalinter",
		"runs a whole bunch of Go linters",
	}, nil}
}

// Setup linters
func (le *lintMessagesEvaluator) Setup() {
	l := make(map[string]*linter)
	l = map[string]*linter{
		"gofmt":       {0, 3, 0, "https://golang.org/cmd/gofmt/", "detects if Go code is incorrectly formatted"},
		"goimports":   {0, 2, 0, "https://golang.org/x/tools/cmd/goimports", "finds missing imports"},
		"golint":      {7, 1, 0, "https://github.com/golang/lint", "official linter for Go code"},
		"dupl":        {5, 1.5, 0, "https://github.com/mibk/dupl", "examines Go code and finds duplicated code"},
		"deadcode":    {0, 3, 0, "https://golang.org/src/cmd/vet/deadcode.go", "checks for syntactically unreachable Go code"},
		"gocyclo":     {6, 2, 0, "https://github.com/fzipp/gocyclo", "calculates cyclomatic complexities of functions in Go code"},
		"vet":         {0, 2.5, 0, "https://golang.org/cmd/vet", "examines Go code and reports suspicious constructs"},
		"vetshadow":   {0, 1, 0, "https://golang.org/src/cmd/vet/shadow.go", "examines Go code and reports shadowed variables"},
		"ineffassign": {0, 1, 0, "https://github.com/gordonklaus/ineffassign", "detects ineffective assignments in Go code"},
		"errcheck":    {6, 2, 0, "https://github.com/kisielk/errcheck", "finds unchecked errors in Go code"},
		"goconst":     {1.5, 1, 0, "https://github.com/jgautheron/goconst", "finds repeated strings in Go code that could be replaced by a constant"},
		"gosimple":    {0, 1.5, 0, "https://github.com/dominikh/go-simple", "examines Go code and reports constructs that can be simplified"},
		"staticcheck": {0, 1.5, 0, "https://github.com/dominikh/go-staticcheck", "checks the inputs to certain functions, such as regexp"},
	}

	le.linters = l
}

// Calculate overloads Evaluator/Calculate
func (le *lintMessagesEvaluator) Calculate(p map[string]interface{}) *model.EvaluatorResponse {
	r := le.NewResponse(100, 2, "", nil)
	lm := p[model.LintMessagesName].(model.LintMessages)
	cs := p[model.CodeStatsName].(model.CodeStats)

	// Loop over messages
	for _, m := range lm {
		for ln, lr := range m {
			if l, ok := le.linters[ln]; ok {
				l.warnings = 0
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
	details := []*model.EvaluatorResponse{}

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

			details = append(details, &model.EvaluatorResponse{
				n,
				score,
				d.weight,
				d.desc,
				"exceeds the warnings/LOC threshold",
				d.url,
				nil,
			})

			log.WithFields(log.Fields{
				"score":  score,
				"weight": d.weight,
			}).Debugf("[%s] score per linter", n)

			scores = append(scores, score*d.weight)
		}
	}

	if len(details) > 0 {
		r.Details = details
	}

	if len(scores) > 0 {
		r.Score = xmath.Sum(scores) / weights
	}

	return r
}
