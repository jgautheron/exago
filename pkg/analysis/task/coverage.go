package task

import (
	"os"
	"time"

	"github.com/jgautheron/exago/pkg/analysis/cov"
)

type coverageRunner struct {
	Runner
	tempFile *os.File
}

// CoverageRunner is a runner used for testing Go projects
func CoverageRunner(m *Manager) Runnable {
	return &coverageRunner{
		Runner: Runner{Label: "Code Coverage", Mgr: m},
	}
}

// Execute gets all the coverage files and returns the output
func (r *coverageRunner) Execute() error {
	defer r.trackTime(time.Now())
	rep, err := cov.ConvertRepository(r.Manager().Repository())
	if err != nil {
		return err
	}

	r.Data = rep

	return nil
}
