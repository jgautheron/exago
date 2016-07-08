package score

import (
	"math"

	"github.com/SimonWaldherr/golibs/xmath"
	log "github.com/Sirupsen/logrus"

	"github.com/exago/svc/repository/model"
)

type linter struct {
	threshold float64
	weight    float64
	drop      float64
	url       string
	desc      string
	warnings  float64
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

	// Linter map, arguments: threshold, weight, drop rate, url and description, last arg is warning counter
	l = map[string]*linter{
		"gofmt":       {0, 3, -1.8, "https://golang.org/cmd/gofmt/", "detects if Go code is incorrectly formatted", 0},
		"goimports":   {0, 2, -0.6, "https://golang.org/x/tools/cmd/goimports", "finds missing imports", 0},
		"golint":      {4, 1, -0.2, "https://github.com/golang/lint", "official linter for Go code", 0},
		"dupl":        {2, 1.5, -0.2, "https://github.com/mibk/dupl", "examines Go code and finds duplicated code", 0},
		"deadcode":    {0, 3, -0.7, "https://golang.org/src/cmd/vet/deadcode.go", "checks for syntactically unreachable Go code", 0},
		"gocyclo":     {3, 2, -0.5, "https://github.com/fzipp/gocyclo", "calculates cyclomatic complexities of functions in Go code", 0},
		"vet":         {0, 2.5, -0.8, "https://golang.org/cmd/vet", "examines Go code and reports suspicious constructs", 0},
		"vetshadow":   {0, 1, -0.7, "https://golang.org/src/cmd/vet/shadow.go", "examines Go code and reports shadowed variables", 0},
		"ineffassign": {0, 1, -0.6, "https://github.com/gordonklaus/ineffassign", "detects ineffective assignments in Go code", 0},
		"errcheck":    {0, 2, -0.2, "https://github.com/kisielk/errcheck", "finds unchecked errors in Go code", 0},
		"goconst":     {1.5, 1, -0.2, "https://github.com/jgautheron/goconst", "finds repeated strings in Go code that could be replaced by a constant", 0},
		"gosimple":    {0, 1.5, -0.3, "https://github.com/dominikh/go-simple", "examines Go code and reports constructs that can be simplified", 0},
		"staticcheck": {0, 1.5, -0.4, "https://github.com/dominikh/go-staticcheck", "checks the inputs to certain functions, such as regexp", 0},
	}

	le.linters = l
}

// Calculate overloads Evaluator/Calculate
func (le *lintMessagesEvaluator) Calculate(p map[string]interface{}) *model.EvaluatorResponse {
	r := le.NewResponse(100, 2, "", nil)
	lm := p[model.LintMessagesName].(model.LintMessages)
	cs := p[model.CodeStatsName].(model.CodeStats)

	// Loop over messages, counting all warnings
	// @todo improve incoming structure so we avoid these ugly nested loops
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
		// Compute the ratio warnings/LOC that we multiply by 100
		tmp := 100 * d.warnings / float64(cs["LOC"])

		log.WithFields(log.Fields{
			"defect ratio": tmp,
			"threshold":    d.threshold,
			"warnings":     d.warnings,
			"loc":          cs["LOC"],
			"weight":       d.weight,
		}).Debugf("[%s] threshold vs ratio", n)

		// If ratio exceeds threshold, calculate linter score
		if tmp > d.threshold {
			// We compute a simple exponential decay based on linter rate decay
			// 100 * exp(drop*ratio)
			score := 100 * math.Exp(d.drop*tmp)
			weights += d.weight

			// Create an evaluator response specific to each linter
			details = append(details, &model.EvaluatorResponse{
				Name:    n,
				Score:   score,
				Weight:  d.weight,
				Desc:    d.desc,
				Message: "exceeds the warnings/LOC threshold",
				URL:     d.url,
				Details: nil,
			})

			log.WithFields(log.Fields{
				"score":  score,
				"weight": d.weight,
			}).Debugf("[%s] score per linter", n)

			scores = append(scores, score*d.weight)
		}
	}

	// If we have details append them to response
	if len(details) > 0 {
		r.Details = details
	}

	// If we have linter scores, compute the weighted average
	if len(scores) > 0 {
		r.Score = xmath.Sum(scores) / weights
	}

	return r
}
