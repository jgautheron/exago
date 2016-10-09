package processor

import (
	"errors"
	"fmt"
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
	logger = log.WithField("prefix", "processor")

	// DefaultTypes represents the default processors enabled.
	DefaultTypes = []string{
		"codestats",
		"projectrunner",
		"lintmessages",
	}
	ErrRoutineTimeout = errors.New("The analysis timed out")
)

type processingError struct {
	tp, message, output string
}

func (e processingError) Error() string {
	return fmt.Sprintf(
		`%s returned the error: "%s"; output: %s`,
		e.tp, e.message, e.output,
	)
}

type errorOutput struct {
	Type   string `json:"type,omitempty"`
	Error  string `json:"error"`
	Output string `json:"output,omitempty"`
}

func ProcessRepository(repo string, tr taskrunner.TaskRunner) (interface{}, error) {
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
	logger         *log.Entry
	types, linters []string
	taskrunner     taskrunner.TaskRunner
	aborted        bool
	processed      chan bool
	Repository     repository.Record
	HasError       bool
	Aborted        chan bool
	Done           chan bool
	Output         map[string]interface{}
}

func NewChecker(repo string, tr taskrunner.TaskRunner) *Checker {
	return &Checker{
		logger:     logger.WithField("repository", repo),
		types:      DefaultTypes,
		linters:    repository.DefaultLinters,
		taskrunner: tr,
		processed:  make(chan bool),
		Repository: repository.New(repo, ""),
		HasError:   false,
		Aborted:    make(chan bool, 1),
		Done:       make(chan bool, 1),
	}
}

// Run launches concurrently every check and merges the output.
func (rc *Checker) Run() {
	rc.Repository.SetStartTime(time.Now())

	for _, tp := range rc.types {
		go func(tp string) {
			var (
				out interface{}
				err error
			)

			switch tp {
			case model.CodeStatsName:
				out, err = rc.taskrunner.FetchCodeStats()
				if err == nil {
					rc.Repository.SetCodeStats(out.(model.CodeStats))
				} else {
					rc.Repository.SetError(tp, err)
				}
			case model.ProjectRunnerName:
				out, err = rc.taskrunner.FetchProjectRunner()
				if err == nil {
					rc.Repository.SetProjectRunner(out.(model.ProjectRunner))
				} else {
					rc.Repository.SetError(tp, err)
				}
			case model.LintMessagesName:
				out, err = rc.taskrunner.FetchLintMessages(rc.linters)
				if err == nil {
					rc.Repository.SetLintMessages(out.(model.LintMessages))
				} else {
					rc.Repository.SetError(tp, err)
				}
			}

			if err != nil {
				rc.HasError = true
				rc.logger.Error(err)
			}

			rc.processed <- true
		}(tp)

		lgr := rc.logger.WithField("type", tp)

		select {
		case <-rc.processed:
			// item processed
		case <-rc.Aborted:
			lgr.Warn("Shutting down (aborted)")
		case <-time.After(RoutineTimeout):
			rc.Repository.SetError(tp, ErrRoutineTimeout)
			lgr.Error(ErrRoutineTimeout)
		}
	}

	if !rc.aborted {
		rc.StampEntry()
		rc.Done <- true
	}
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
	rc.aborted = true
	close(rc.Aborted)
}
