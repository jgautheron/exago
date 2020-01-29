package task

import (
	"encoding/json"
	"os"
	"os/exec"
	"time"

	exago "github.com/jgautheron/exago/pkg"

	"github.com/pkg/errors"
)

type lintRunner struct {
	Runner
}

type LinterResponse struct {
	Issues []LinterIssue
}

type LinterIssue struct {
	FromLinter  string
	Text        string
	SourceLines []string
	Replacement string
	Pos         LinterIssuePos
}

type LinterIssuePos struct {
	Filename string
	Offset   int
	Line     int
	Column   int
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
	p := []string{"run", "--out-format=json", "--issues-exit-code=0"}
	rep := r.Manager().RepositoryPath()
	if r.Manager().Reference() != "" {
		rep += ":" + r.Manager().Reference()
	}
	p = append(p, rep+"/...")

	os.Setenv("GO111MODULE", "off")
	out, err := exec.Command("golangci-lint", p...).CombinedOutput()

	if err != nil {
		// If we cannot run linter return with error
		return errors.Wrap(err, string(out))
	}

	var linterOutput LinterResponse
	linterResults := exago.LinterResults{}

	json.Unmarshal(out, &linterOutput)

	// Format to something like:
	//  example_file.go
	//	  linterX: messages / issues
	//    linterY: messages / issues
	for _, issue := range linterOutput.Issues {
		if _, ok := linterResults[issue.Pos.Filename]; !ok {
			// Means file was empty
			linterResults[issue.Pos.Filename] = []exago.LinterResult{
				{Linter: issue.FromLinter, Messages: []exago.LinterMessage{
					{Column: issue.Pos.Column, Message: issue.Text, Row: issue.Pos.Line, Severity: "error"},
				}}}
		} else {
			// Means there is a file already
			// we have to check if there is already a linter for that file
			linterIdx := -1
			for idx, linterFile := range linterResults[issue.Pos.Filename] {
				if linterFile.Linter == issue.FromLinter {
					linterIdx = idx
					break
				}
			}

			if linterIdx >= 0 {
				// Means we have given linter for current file
				// We have to append only new message
				linterResults[issue.Pos.Filename][linterIdx].Messages = append(linterResults[issue.Pos.Filename][linterIdx].Messages, exago.LinterMessage{
					Column:   issue.Pos.Column,
					Message:  issue.Text,
					Row:      issue.Pos.Line,
					Severity: "error",
				})
			} else {
				// This linter doesn't exist yet for current file, we have to add it
				linterResults[issue.Pos.Filename] = append(linterResults[issue.Pos.Filename], exago.LinterResult{
					Linter:   issue.FromLinter,
					Messages: []exago.LinterMessage{{Column: issue.Pos.Column, Message: issue.Text, Row: issue.Pos.Line, Severity: "error"}},
				})
			}

		}
	}

	r.Data = linterResults
	return nil
}
