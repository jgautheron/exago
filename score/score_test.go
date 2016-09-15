package score_test

import (
	"testing"

	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/score"
)

func TestScore(t *testing.T) {
	d := getStubData(2500, 200, 0.8, 75, 5, []string{"projectBuilds", "isFormatted", "hasReadme", "isDirMatch"})
	sc, _ := score.Process(d)
	if sc < 80 {
		t.Error("The score should exceed 80 pts")
	}
}

func getStubData(loc int, cloc int, duration, coverage float64, thirdParties int, checklist []string) model.Data {
	d := model.Data{}

	p := []model.Package{}
	p = append(p, model.Package{ExecutionTime: duration, Coverage: coverage})

	d.CodeStats = map[string]int{"LOC": loc, "CLOC": cloc, "Test": 123}
	d.ProjectRunner = model.ProjectRunner{
		Packages:     p,
		ThirdParties: getThirdParties(thirdParties),
		Checklist:    getStubChecklist(checklist),
	}
	d.LintMessages = getStubMessages(map[string]int{"gas": 3})

	return d
}
