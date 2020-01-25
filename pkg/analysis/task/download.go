package task

import (
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

type downloadRunner struct {
	Runner
}

// DownloadRunner is a runner used for downloading Go projects
// from remote repositories such as Github, Bitbucket etc.
func DownloadRunner(m *Manager) Runnable {
	return &downloadRunner{
		Runner: Runner{Label: "Go Get", Mgr: m},
	}
}

// Execute, downloads a Go repository using the go get command
// too bad, we can't do this as a library :/
func (r *downloadRunner) Execute() error {
	defer r.trackTime(time.Now())

	// Return early if repository is already in the GOPATH
	if _, err := os.Stat(r.Manager().RepositoryPath()); err == nil {
		return r.toRepoDir()
	}

	// Go get the package
	p := []string{"get", "-d", "-t"}
	rep := r.Manager().Repository()
	if r.Manager().Reference() != "" {
		rep += ":" + r.Manager().Reference()
	}
	p = append(p, rep+"/...")

	os.Setenv("GO111MODULE", "off")
	out, err := exec.Command("go", p...).CombinedOutput()
	if err != nil {
		// If we can't download, stop execution as BreakOnError is true with this runner
		return errors.Wrap(err, string(out))
	}

	r.RawOutput = string(out)

	// cd into repository
	return r.toRepoDir()
}

func (r *downloadRunner) toRepoDir() error {
	// Change directory
	err := os.Chdir(r.Manager().RepositoryPath())
	if err != nil {
		return err
	}

	return nil
}
