package task

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Manager contains all registered runnables
type Manager struct {
	Success bool                `json:"success"`
	Errors  map[string]string   `json:"errors,omitempty"`
	Runners map[string]Runnable `json:"data,omitempty"`

	repository     string
	repositoryPath string
	reference      string
}

// NewManager instantiates a runnable manager
// the manager has the responsibility to execute all runners
// and decide whether a runner should run in parallel processing or not
func NewManager(r string) *Manager {
	m := &Manager{
		repository:     r,
		repositoryPath: fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), r),
		Errors:         make(map[string]string),
	}

	if strings.TrimSpace(r) == "" {
		m.Errors[downloadName] = "Repository is required"
		return m
	}

	m.Runners = map[string]Runnable{
		downloadName:     DownloadRunner(m),
		locName:          LocRunner(m),
		testName:         TestRunner(m),
		coverageName:     CoverageRunner(m),
		checklistName:    ChecklistRunner(m),
		thirdPartiesName: ThirdPartiesRunner(m),
	}

	return m
}

// UseReference sets reference flag
func (m *Manager) UseReference(r string) {
	m.reference = r
}

// Reference returns reference
func (m *Manager) Reference() string {
	return m.reference
}

// RepositoryPath returns repository path
func (m *Manager) RepositoryPath() string {
	return m.repositoryPath
}

// Repository returns repository (e.g. :vcs/:owner/:package+)
func (m *Manager) Repository() string {
	return m.repository
}

// ExecuteRunners launches the runners
func (m *Manager) ExecuteRunners() Manager {
	// Execute download runner synchronously
	dlr := m.Runners[downloadName]
	// Execute synchronously
	err := dlr.Execute()
	// Exit early if we can't download
	if err != nil {
		m.Errors[downloadName] = err.Error()
		m.Runners = nil
		return m
	}

	var wg sync.WaitGroup
	for n, ru := range m.Runners {
		// Skip download runner
		if n == downloadName {
			continue
		}
		// Increment the WaitGroup counter.
		wg.Add(1)
		go func(r Runnable, name string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Execute the runner
			err := r.Execute()
			if err != nil {
				m.Errors[name] = err.Error()
			}
		}(ru, n)
	}

	// Wait for all runners to complete.
	wg.Wait()

	if len(m.Errors) == 0 {
		m.Success = true
	} else {
		m.Runners = nil
	}

	return m
}
