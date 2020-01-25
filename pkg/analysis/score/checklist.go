package score

import (
	exago "github.com/jgautheron/exago/pkg"
	"github.com/sirupsen/logrus"

	"simonwaldherr.de/go/golibs/xmath"
)

type checker struct {
	score  float64
	weight float64
	url    string
}

type checkListEvaluator struct {
	Evaluator
	checkers map[string]*checker
}

// CheckListEvaluator measures a score based on given checklist criterias
func CheckListEvaluator() CriteriaEvaluator {
	return &checkListEvaluator{Evaluator{
		"checklist",
		"https://github.com/jgautheron/exago",
		"inspects project for best practices",
	}, nil}
}

// Setup checkers
func (ce *checkListEvaluator) Setup() {
	c := make(map[string]*checker)
	c = map[string]*checker{
		"isFormatted":     {100, 3, "https://github.com/matttproud/gochecklist/blob/master/publication/code_correctness.md"},
		"isDirMatch":      {100, .7, "https://github.com/matttproud/gochecklist/blob/master/publication/dir_pkg_match.md"},
		"isLinted":        {100, .5, "https://github.com/matttproud/gochecklist/blob/master/publication/code_correctness.md"},
		"isVetted":        {100, .5, "https://github.com/matttproud/gochecklist/blob/master/publication/govet_correctness.md"},
		"hasReadme":       {100, 3, "https://github.com/matttproud/gochecklist/blob/master/publication/documentation_entrypoint.md"},
		"hasBenches":      {100, .5, ""},
		"hasContributing": {100, .3, ""},
	}

	ce.checkers = c
}

// Calculate overloads Evaluator/Calculate
func (ce *checkListEvaluator) Calculate(d exago.Data) *exago.EvaluatorResponse {
	r := ce.NewResponse(100, 1.8, "", nil)
	cl := d.Results

	for _, failed := range cl.Checklist.Data.Failed {
		if ch, ok := ce.checkers[failed]; ok {
			ch.score = 0
		}
	}

	// Compute score
	scores := []float64{}
	weights := 0.0
	details := []*exago.EvaluatorResponse{}

	for n, c := range ce.checkers {
		weights += c.weight
		msg := "check failed"
		if c.score == 100 {
			msg = "check succeeded"
		}
		details = append(details, &exago.EvaluatorResponse{
			n,
			c.score,
			c.weight,
			"",
			msg,
			c.url,
			nil,
		})

		logrus.WithFields(logrus.Fields{
			"score":  c.score,
			"weight": c.weight,
		}).Debugf("[%s] score per checker", n)

		scores = append(scores, c.score*c.weight)
	}

	if len(details) > 0 {
		r.Details = details
	}

	if len(scores) > 0 {
		r.Score = xmath.Sum(scores) / weights
	}

	return r
}
