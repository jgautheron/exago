package task

import (
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type thirdPartiesRunner struct {
	Runner
}

// ThirdPartiesRunner launches go list to find all dependencies
func ThirdPartiesRunner(m *Manager) Runnable {
	return &thirdPartiesRunner{
		Runner{Label: "Go List (finds all 3rd parties)", Mgr: m},
	}
}

// Execute go list
func (r *thirdPartiesRunner) Execute() error {
	defer r.trackTime(time.Now())

	list, err := exec.Command("go", "list", "-f", `'{{ join .Deps ", " }}'`, "./...").CombinedOutput()
	if err != nil {
		return err
	}

	r.Data = r.parseListOutput(string(list))

	return nil
}

func (r *thirdPartiesRunner) parseListOutput(output string) (out []string) {
	reg := regexp.MustCompile(`(?m)([\w\d\-]+)\.([\w]{2,})\/([\w\d\-]+)\/([\w\d\-\.]+)(\.v\d+)?`)
	out = make([]string, 0)
	uniq := map[string]bool{}
	sl := strings.Split(output, ",")
	for _, v := range sl {
		v = strings.TrimSpace(v)
		m := reg.FindAllString(v, -1)

		// Only interested in third parties
		if len(m) == 0 {
			continue
		}

		// Match only the last path found in the path
		// That way we support imports made this way:
		// github.com/heroku/hk/Godeps/_workspace/src/code.google.com/p/go-uuid/uuid
		lastMatch := m[len(m)-1]
		if lastMatch == r.Manager().Repository() {
			continue
		}

		uniq[lastMatch] = true
	}

	for pkg := range uniq {
		out = append(out, pkg)
	}

	return out
}
