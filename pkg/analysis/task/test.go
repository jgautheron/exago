package task

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type testRunner struct {
	Runner
}

type pkg struct {
	Name          string  `json:"name,omitempty"`
	ExecutionTime float64 `json:"execution_time,omitempty"`
	Success       bool    `json:"success"`
	Tests         []test  `json:"tests"`
}

type test struct {
	Name          string  `json:"name,omitempty"`
	ExecutionTime float64 `json:"execution_time"`
	Passed        bool    `json:"passed"`
}

// TestRunner is a runner used for testing Go projects
func TestRunner(m *Manager) Runnable {
	return &testRunner{
		Runner{Label: "Go Test", Mgr: m},
	}
}

// Execute tests and determine which tests are passing/failing
func (r *testRunner) Execute() error {
	defer r.trackTime(time.Now())

	out, err := exec.Command("bash", "-c", "go test -v $(go list ./... | grep -v vendor | grep -v Godeps)").CombinedOutput()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return errors.New(string(out))
		}
	}

	r.RawOutput = string(out)
	r.parseTestOutput()

	return nil
}

// parseTestOutput parses test output and fills Data property
func (r *testRunner) parseTestOutput() {
	pkgs, p, tests := []pkg{}, pkg{}, []test{}

	pkgReregex, _ := regexp.Compile(`(?i)^(ok|fail|\?)\s+([\w\.\/\-]+)(?:\s+([\d\.]+)s)?(?:\s+coverage:\s(\d+(\.\d+)?))?`)
	testRegex, _ := regexp.Compile(`(?i)(FAIL|PASS):\s([\w\d]+)\s\(([\d\.]+)s\)`)

	lines := strings.Split(r.RawOutput, "\n")
	for _, l := range lines {
		submatch := testRegex.FindStringSubmatch(l)
		if len(submatch) != 0 {
			ef, _ := strconv.ParseFloat(submatch[3], 64)
			tests = append(tests, test{
				submatch[2],
				ef,
				submatch[1] == "PASS",
			})
			continue
		}

		submatch = pkgReregex.FindStringSubmatch(l)
		if len(submatch) != 0 {
			switch submatch[1] {
			case "FAIL":
				ef, _ := strconv.ParseFloat(submatch[3], 64)
				p = pkg{
					Name:          submatch[2],
					Success:       false,
					Tests:         tests,
					ExecutionTime: ef,
				}
			case "ok":
				ef, _ := strconv.ParseFloat(submatch[3], 64)
				p = pkg{
					Name:          submatch[2],
					Success:       true,
					Tests:         tests,
					ExecutionTime: ef,
				}
			default:
				continue
			}

			pkgs = append(pkgs, p)
			tests = []test{}
		}
	}

	r.Data = pkgs
}
