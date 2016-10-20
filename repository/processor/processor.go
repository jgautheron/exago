package processor

import (
	"errors"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/taskrunner"
)

const (
	// Lambda function time limit
	RoutineTimeout = time.Second * 280
)

var (
	logger            = log.WithField("prefix", "processor")
	ErrRoutineTimeout = errors.New("The analysis timed out")
)

func ProcessRepository(repo string, tr taskrunner.TaskRunner) (interface{}, error) {
	if _, err := repository.IsValid(repo); err != nil {
		return nil, err
	}

	checker := NewChecker(repo, tr)
	checker.Run()

	var out model.Data
	select {
	case <-checker.Done:
		out = checker.Repository.GetData()
	case <-checker.Aborted:
		out = checker.Repository.GetData()
	}
	return out, nil
}

type Checker struct {
	logger     *log.Entry
	taskrunner taskrunner.TaskRunner
	Repository repository.Record
	Aborted    chan bool
	Done       chan bool
}

func NewChecker(repo string, tr taskrunner.TaskRunner) *Checker {
	return &Checker{
		logger:     logger.WithField("repository", repo),
		taskrunner: tr,
		Repository: repository.New(repo, ""),
		Aborted:    make(chan bool, 1),
		Done:       make(chan bool, 1),
	}
}

// Run launches concurrently every check and merges the output.
func (rc *Checker) Run() {
	rc.Repository.SetStartTime(time.Now())

	var out interface{}
	var err error

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func(rc *Checker) {
		defer wg.Done()
		out, err = rc.taskrunner.FetchCodeStats()
		if err == nil {
			rc.Repository.SetCodeStats(out.(model.CodeStats))
		} else {
			rc.Repository.SetError(model.CodeStatsName, err)
		}
	}(rc)

	go func(rc *Checker) {
		defer wg.Done()
		out, err = rc.taskrunner.FetchProjectRunner()
		if err == nil {
			rc.Repository.SetProjectRunner(out.(model.ProjectRunner))
		} else {
			rc.Repository.SetError(model.ProjectRunnerName, err)
		}
	}(rc)

	go func(rc *Checker) {
		defer wg.Done()
		out, err = rc.taskrunner.FetchLintMessages()
		if err == nil {
			rc.Repository.SetLintMessages(out.(model.LintMessages))
		} else {
			rc.Repository.SetError(model.LintMessagesName, err)
		}
	}(rc)

	wg.Wait()
	rc.StampEntry()
	rc.Done <- true
}

// StampEntry is called once the entire dataset is loaded.
func (rc *Checker) StampEntry() {
	// Add the metadata
	err := rc.Repository.SetMetadata()
	if err != nil {
		rc.Repository.SetError(model.MetadataName, err)
	}

	// Add the score
	err = rc.Repository.SetScore()
	if err != nil {
		rc.Repository.SetError(model.ScoreName, err)
	}

	// Add the timestamp
	rc.Repository.SetLastUpdate()

	// Add the execution time
	rc.Repository.SetExecutionTime()

	// Persist the dataset
	if err := rc.Repository.Save(); err != nil {
		rc.logger.Errorf("Could not persist the dataset: %v", err)
	}
}

// Abort declares the task as done and skips the processing.
func (rc *Checker) Abort() {
	close(rc.Aborted)
}
