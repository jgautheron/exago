package processor

import (
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/repository"
	"github.com/exago/svc/repository/model"
	"github.com/exago/svc/showcaser"
	"github.com/exago/svc/taskrunner"
)

const (
	// Lambda function time limit
	RoutineTimeout = time.Minute * 5
)

var (
	ErrRoutineTimeout = errors.New("The analysis timed out")

	// DefaultTypes represents the default processors enabled.
	DefaultTypes = []string{
		"imports",
		"codestats",
		"testresults",
		"lintmessages",
	}
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

type Checker struct {
	logger         *log.Entry
	types, linters []string
	taskrunner     taskrunner.TaskRunner
	processed      chan bool
	Repository     repository.Record
	HasError       bool
	Aborted        chan bool
	Done           chan bool
	Output         map[string]interface{}
}

func NewChecker(repo string, tr taskrunner.TaskRunner) *Checker {
	return &Checker{
		logger:     log.WithField("repository", repo),
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
			case model.ImportsName:
				out, err = rc.taskrunner.FetchImports()
				if err == nil {
					rc.Repository.SetImports(out.(model.Imports))
				} else {
					rc.Repository.SetError(tp, err)
				}
			case model.CodeStatsName:
				out, err = rc.taskrunner.FetchCodeStats()
				if err == nil {
					rc.Repository.SetCodeStats(out.(model.CodeStats))
				} else {
					rc.Repository.SetError(tp, err)
				}
			case model.TestResultsName:
				out, err = rc.taskrunner.FetchTestResults()
				if err == nil {
					rc.Repository.SetTestResults(out.(model.TestResults))
				} else {
					// Expose isolated errors
					switch ts := out.(model.TestResults); {
					case ts.Errors.Goget != "":
						err = processingError{"goget", ts.Errors.Goget, ts.RawOutput.Goget}
					case ts.Errors.Gotest != "":
						err = processingError{"gotest", ts.Errors.Gotest, ts.RawOutput.Gotest}
					}
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
				log.Error(err)
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
			rc.Repository.SetError("", ErrRoutineTimeout)
			lgr.Error(ErrRoutineTimeout)
		}
	}

	rc.StampEntry()
	rc.Done <- true
	go showcaser.ProcessRepository(rc.Repository)
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
		log.Errorf("Could not persist the dataset: %v", err)
	}
}

// Abort declares the task as done and skips the processing.
func (rc *Checker) Abort() {
	rc.Aborted <- true
}
