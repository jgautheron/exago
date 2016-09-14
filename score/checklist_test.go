package score_test

import (
	"testing"

	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/score"
)

var criterias = []string{"projectBuilds", "isFormatted", "hasReadme", "isDirMatch", "isLinted", "hasBenches"}

func TestNotEvenOne(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = getStubChecklist([]string{})
	d.CodeStats = map[string]int{"LOC": 200}
	evaluator := score.CheckListEvaluator()
	evaluator.Setup()
	res := evaluator.Calculate(d)
	if res.Score != 0 {
		t.Error("The score should be 0")
	}
}

func TestFmtFail(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = getStubChecklist([]string{"projectBuilds", "hasReadme", "isDirMatch", "isLinted", "hasBenches"})
	d.CodeStats = map[string]int{"LOC": 200}
	evaluator := score.CheckListEvaluator()
	evaluator.Setup()
	res := evaluator.Calculate(d)

	// If a project is not gofmt'd, it probably means we're dealing with a beginner or old project
	if res.Score > 70 {
		t.Error("The score should not exceed 70")
	}
}

func TestReadmeFail(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = getStubChecklist([]string{"projectBuilds", "isFormatted", "isDirMatch", "isLinted", "hasBenches"})
	d.CodeStats = map[string]int{"LOC": 200}
	evaluator := score.CheckListEvaluator()
	evaluator.Setup()
	res := evaluator.Calculate(d)

	// If a project doesn't have a README, it usually means it hasn't been finished yet
	// Creating a README is usually the first step toward open-sourcing a project
	if res.Score > 70 {
		t.Error("The score should not exceed 70")
	}
}

func getStubChecklist(passed []string) model.ProjectRunner {
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

	failedItemList := []model.ChecklistItem{}
	for _, item := range failed {
		failedItemList = append(failedItemList, model.ChecklistItem{
			Name: item,
		})
	}

	passedItemList := []model.ChecklistItem{}
	for _, item := range passed {
		passedItemList = append(passedItemList, model.ChecklistItem{
			Name: item,
		})
	}

	return model.ProjectRunner{
		Checklist: struct {
			Failed []model.ChecklistItem `json:"Failed"`
			Passed []model.ChecklistItem `json:"Passed"`
		}{
			Failed: failedItemList,
			Passed: passedItemList,
		},
	}
}
