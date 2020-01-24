package task

import (
	"time"

	"github.com/jgautheron/golocc"
)

const ignore = `vendor|Godeps|external|pb\.go|bindata\.go|yacc|mocks`

type locRunner struct {
	Runner
}

// LocRunner is a runner used for counting lines of code, comments
// functions, structs, imports etc.
func LocRunner(m *Manager) Runnable {
	return &locRunner{
		Runner: Runner{Label: "Lines of code: LOC, CLOC, NCLOC", Mgr: m},
	}
}

// Execute calls the golocc library
func (r *locRunner) Execute() error {
	defer r.trackTime(time.Now())

	parser := golocc.New(r.Manager().RepositoryPath(), ignore, true)
	res, err := parser.ParseTree()
	if err != nil {
		return err
	}

	r.Data = res

	return nil
}
