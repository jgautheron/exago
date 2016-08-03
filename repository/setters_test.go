package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/exago/svc/repository/model"
)

func TestImportsChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	im := model.Imports([]string{"foo", "bar", "moo"})
	rp.SetImports(im)
	if len(rp.GetImports()) != len(im) {
		t.Error("The imports have not changed")
	}
}

func TestCodeStatsChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	cs := rp.GetCodeStats()
	cs["CLOC"] = 123
	rp.SetCodeStats(cs)
	if rp.GetCodeStats()["CLOC"] != cs["CLOC"] {
		t.Error("The CLOC has not changed")
	}
}

func TestLintMessagesChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	lm := rp.GetLintMessages()
	lm["codename.go"]["errcheck"][0]["col"] = 123
	rp.SetLintMessages(lm)
	if rp.GetLintMessages()["codename.go"]["errcheck"][0]["col"] != lm["codename.go"]["errcheck"][0]["col"] {
		t.Error("The col has not changed")
	}
}

func TestStartTimeChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	now := time.Now()
	rp.SetStartTime(now)
	if rp.startTime != now {
		t.Error("The start time has not changed")
	}
}

func TestExecutionTimeChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	now := time.Now()
	rp.SetStartTime(now)
	rp.SetExecutionTime()
	if rp.GetExecutionTime() != "0s" {
		t.Error("The execution time has not changed")
	}
}

func TestLastUpdateTimeChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	now := time.Now()
	rp.SetLastUpdate()
	if rp.GetLastUpdate().Day() != now.Day() {
		t.Error("The last update time has not changed")
	}
}

func TestScoreChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	rp.SetScore()
	if rp.GetScore().Rank != "D" {
		t.Error("The rank has not changed")
	}
}

func TestErrorAdded(t *testing.T) {
	rp, _ := loadStubRepo()
	rp.SetError("codestats", errors.New("Could not load code stats!"))
	if len(rp.Data.Errors) != 1 {
		t.Error("The error has not been added")
	}
}
