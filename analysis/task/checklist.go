package task

import (
	"time"

	"github.com/jgautheron/exago/analysis/checklist"
)

type checklistRunner struct {
	Runner
}

// ChecklistRunner runs the checklist
func ChecklistRunner(m *Manager) Runnable {
	return &checklistRunner{
		Runner{Label: "CheckList", Mgr: m},
	}
}

// Execute checklist
func (r *checklistRunner) Execute() error {
	defer r.trackTime(time.Now())

	cl := checklist.New(r.Manager().RepositoryPath())
	passed, failed := cl.RunTasks()

	r.Data = map[string][]string{
		"passed": passed,
		"failed": failed,
	}

	return nil
}
