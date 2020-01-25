package score_test

import (
	"testing"

	exago "github.com/jgautheron/exago/pkg"

	"github.com/jgautheron/exago/pkg/analysis/score"
	"github.com/sirupsen/logrus"
)

func TestScore(t *testing.T) {
	d := getStubData(2500, 200, 0.8, 75, 5, []string{"projectBuilds", "isFormatted", "hasReadme", "isDirMatch"})
	sc, _ := score.Process(d)
	if sc < 80 {
		logrus.Warnln(sc)
		t.Error("The score should exceed 80 pts")
	}
}

func getStubData(loc int, cloc int, duration, coverage float64, thirdParties int, checklist []string) exago.Data {
	d := exago.Data{}

	pr := exago.Results{}
	pr.Coverage.Data.Coverage = coverage
	pr.Thirdparties.Data = getThirdParties(thirdParties)
	pr.Goprove.Data = getStubChecklist(checklist)
	pr.CodeStats.Data = map[string]int{"loc": loc, "cloc": cloc, "test": 123}
	d.ProjectRunner = pr

	d.LintMessages = getStubMessages(map[string]int{"gas": 3})

	return d
}
