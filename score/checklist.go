package score

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/repository/model"

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
		model.CheckListName,
		"https://github.com/karolgorecki/goprove",
		"inspects project for the best practices listed in the Go CheckList",
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
func (ce *checkListEvaluator) Calculate(d model.Data) *model.EvaluatorResponse {
	r := ce.NewResponse(100, 1.8, "", nil)
	cl := d.ProjectRunner

	for _, failed := range cl.Goprove.Data.Failed {
		if ch, ok := ce.checkers[failed.Name]; ok {
			ch.score = 0
		}
	}

	// Compute score
	scores := []float64{}
	weights := 0.0
	details := []*model.EvaluatorResponse{}

	for n, c := range ce.checkers {
		weights += c.weight
		msg := "check failed"
		if c.score == 100 {
			msg = "check succeeded"
		}
		details = append(details, &model.EvaluatorResponse{
			n,
			c.score,
			c.weight,
			"",
			msg,
			c.url,
			nil,
		})

		logger.WithFields(log.Fields{
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
