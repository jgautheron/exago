package score_test

import (
	"testing"

	exago "github.com/jgautheron/exago/pkg"

	"github.com/jgautheron/exago/analysis/score"
)

var criterias = []string{"projectBuilds", "isFormatted", "hasReadme", "isDirMatch", "isLinted", "isVetted", "hasContributing", "hasBenches"}

func TestChecklist(t *testing.T) {
	var tests = []struct {
		criterias []string
		loc       int
		operator  string
		expected  float64
		desc      string
	}{
		{[]string{}, 200, "=", 0, "The score should be 0"},
		{
			[]string{"projectBuilds", "hasReadme", "isDirMatch", "isLinted", "hasBenches"},
			5000, "<", 70,
			"If a project is not gofmt'd, it probably means we're dealing with a beginner or old project",
		},
		{
			[]string{"projectBuilds", "isFormatted", "isDirMatch", "isLinted", "hasBenches"},
			5000, "<", 70,
			"The README file is a documentation entry point, generally a must-have",
		},
	}

	for _, tt := range tests {
		d := exago.Data{}
		d.ProjectRunner = exago.Results{}
		d.ProjectRunner.Goprove.Data = getStubChecklist(tt.criterias)
		d.ProjectRunner.CodeStats.Data = map[string]int{"LOC": 123}
		evaluator := score.CheckListEvaluator()
		evaluator.Setup()
		res := evaluator.Calculate(d)

		switch tt.operator {
		case "<":
			if res.Score > tt.expected {
				t.Errorf("Wrong score %s", tt.desc)
			}
		case ">":
			if res.Score < tt.expected {
				t.Errorf("Wrong score %s", tt.desc)
			}
		case "=":
			if res.Score != tt.expected {
				t.Errorf("Wrong score %s", tt.desc)
			}
		}
	}
}

func getStubChecklist(passed []string) exago.Checklist {
	failed := []string{}
	for _, criteria := range criterias {
		found := false
		for _, item := range passed {
			if item == criteria {
				found = true
			}
		}
		if !found {
			failed = append(failed, criteria)
		}
	}

	failedItemList := []exago.ChecklistItem{}
	for _, item := range failed {
		failedItemList = append(failedItemList, exago.ChecklistItem{
			Name: item,
		})
	}

	passedItemList := []exago.ChecklistItem{}
	for _, item := range passed {
		passedItemList = append(passedItemList, exago.ChecklistItem{
			Name: item,
		})
	}

	return exago.Checklist{
		Failed: failedItemList,
		Passed: passedItemList,
	}
}
