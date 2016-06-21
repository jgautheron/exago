package score

import (
	"log"

	"github.com/exago/svc/repository/model"
)

func init() {
	Register(model.ImportsName, &ImportsEvaluator{Evaluator{100, .5, ""}})
}

type ImportsEvaluator struct {
	Evaluator
}

func (ie *ImportsEvaluator) Calculate(p map[string]interface{}) {
	imp := p[model.ImportsName].(model.Imports)
	tp := len(imp)
	switch true {
	case tp < 4:
		ie.score = 75
		ie.msg = "< 4"
	case tp < 6:
		ie.score = 50
		ie.msg = "< 6"
	case tp < 8:
		ie.score = 25
		ie.msg = "< 8"
	case tp > 8:
		ie.score = 0
		ie.msg = "> 8"
	}

	log.Println(ie)
}
