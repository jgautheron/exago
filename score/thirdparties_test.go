package score_test

import (
	"testing"

	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/score"
)

func TestLowThirdParties(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = model.ProjectRunner{
		ThirdParties: []string{"1", "2"},
	}
	d.CodeStats = map[string]int{"LOC": 200}
	res := score.ThirdPartiesEvaluator().Calculate(d)

	// Two third parties for a small project is pretty common
	if res.Score < 70 || res.Score > 80 {
		t.Error("Wrong score")
	}
}

func TestLotsOfThirdParties(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = model.ProjectRunner{
		ThirdParties: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
	}
	d.CodeStats = map[string]int{"LOC": 5000}
	res := score.ThirdPartiesEvaluator().Calculate(d)

	// 8 third parties for 5000 LOC is not that bad
	if res.Score < 50 {
		t.Error("Wrong score")
	}
}

func TestTooMuchThirdParties(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = model.ProjectRunner{
		ThirdParties: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
	}
	d.CodeStats = map[string]int{"LOC": 2000}
	res := score.ThirdPartiesEvaluator().Calculate(d)

	// For 2k LOC this is proportionally too much
	if res.Score > 50 {
		t.Error("Wrong score")
	}
}

func TestWayTooMuchThirdParties(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = model.ProjectRunner{
		ThirdParties: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
	}
	d.CodeStats = map[string]int{"LOC": 400}
	res := score.ThirdPartiesEvaluator().Calculate(d)

	// For 2k LOC this is proportionally too much
	if res.Score > 30 {
		t.Error("Wrong score")
	}
}

func TestNoThirdParty(t *testing.T) {
	d := model.Data{}
	d.ProjectRunner = model.ProjectRunner{
		ThirdParties: []string{},
	}
	d.CodeStats = map[string]int{"LOC": 100}
	res := score.ThirdPartiesEvaluator().Calculate(d)

	// No third party, then obviously we get the maximum score
	if res.Score != 100 {
		t.Error("The score should be 100")
	}
}
