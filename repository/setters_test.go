package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/hotolab/exago-svc/repository/model"
)

func TestCodeStatsChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	cs := rp.GetCodeStats()
	cs["CLOC"] = 123
	rp.SetCodeStats(cs)
	if rp.GetCodeStats()["CLOC"] != cs["CLOC"] {
		t.Error("The CLOC has not changed")
	}
}

func TestProjectRunnerChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	pr := rp.GetProjectRunner()
	pr.Thirdparties.Data = append(pr.Thirdparties.Data, "github.com/bar/moo")
	rp.SetProjectRunner(pr)
	if len(pr.Thirdparties.Data) != 1 {
		t.Error("The third parties have not changed")
	}
}

func TestLintMessagesChanged(t *testing.T) {
	rp, _ := loadStubRepo()
	lm := rp.GetLintMessages()
	lm["codename.go"]["golint"][0]["col"] = 123
	rp.SetLintMessages(lm)
	if rp.GetLintMessages()["codename.go"]["golint"][0]["col"] != lm["codename.go"]["golint"][0]["col"] {
		t.Error("The col has not changed")
	}
}

func TestStartTimeChanged(t *testing.T) {
	rp := &Repository{
		Name: repo,
	}
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
	rp.SetExecutionTime(time.Since(now))

	et := rp.GetExecutionTime()
	if et != "0s" {
		t.Errorf("The execution time has not changed, got %s", et)
	}
}

func TestLastUpdateTimeChanged(t *testing.T) {
	rp := &Repository{
		Name: repo,
	}
	now := time.Now()
	rp.SetLastUpdate(now)
	if rp.GetLastUpdate() != now {
		t.Error("The last update time has not changed")
	}
}

func TestMetadataChanged(t *testing.T) {
	m := map[string]interface{}{
		"avatar_url":  "http://foo.com/img.png",
		"description": "repository description",
		"language":    "go",
		"stargazers":  123,
		"last_push":   time.Now(),
	}
	rp := &Repository{
		Name: repo,
	}
	rp.SetMetadata(model.Metadata{
		Image:       m["avatar_url"].(string),
		Description: m["description"].(string),
		Stars:       m["stargazers"].(int),
	})
	if rp.GetMetadata().Stars != 123 {
		t.Error("The metadata has not changed")
	}
}

func TestErrorAdded(t *testing.T) {
	rp, _ := loadStubRepo()
	rp.SetError("codestats", errors.New("Could not load code stats!"))
	if len(rp.Data.Errors) != 1 {
		t.Error("The error has not been added")
	}
}
