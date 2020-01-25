package task

import (
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

type lintRunner struct {
	Runner
}

// LintRunner is a runner used for linting files
func LintRunner(m *Manager) Runnable {
	return &lintRunner{
		Runner: Runner{Label: "Go Lint (golangci-lint)", Mgr: m},
	}
}

// Execute runs linters for files using golangci-lint
func (r *lintRunner) Execute() error {
	defer r.trackTime(time.Now())

	// Run linter
	p := []string{"run", "--out-format=json"}
	rep := r.Manager().RepositoryPath()
	if r.Manager().Reference() != "" {
		rep += ":" + r.Manager().Reference()
	}
	p = append(p, rep+"/...")

	os.Setenv("GO111MODULE", "off")
	out, err := exec.Command("golangci-lint", p...).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(out))
	}

	r.RawOutput = string(out)

	return nil
}
