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
	data           chan interface{}
	taskrunner     taskrunner.TaskRunner
	Repository     repository.RepositoryData
	HasError       bool
	Errors         chan error
	Aborted        chan bool
	Done           chan bool
	Output         map[string]interface{}
}

func NewChecker(repo string, tr taskrunner.TaskRunner) *Checker {
	return &Checker{
		logger:     log.WithField("repository", repo),
		types:      DefaultTypes,
		linters:    repository.DefaultLinters,
		data:       make(chan interface{}),
		taskrunner: tr,
		Repository: repository.New(repo, ""),
		HasError:   false,
		Errors:     make(chan error),
		Aborted:    make(chan bool, 1),
		Done:       make(chan bool, 1),
		Output:     map[string]interface{}{},
	}
}

// Run launches concurrently every check and merges the output.
func (rc *Checker) Run() {
	rc.Repository.SetStartTime(time.Now())

	i := 0
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
				}
			case model.CodeStatsName:
				out, err = rc.taskrunner.FetchCodeStats()
				if err == nil {
					rc.Repository.SetCodeStats(out.(model.CodeStats))
				}
			case model.TestResultsName:
				out, err = rc.taskrunner.FetchTestResults()
				if err == nil {
					rc.Repository.SetTestResults(out.(model.TestResults))
				}

				// Expose isolated errors
				switch ts := out.(model.TestResults); {
				case ts.Errors.Goget != "":
					err = processingError{"goget", ts.Errors.Goget, ts.RawOutput.Goget}
				case ts.Errors.Gotest != "":
					err = processingError{"gotest", ts.Errors.Gotest, ts.RawOutput.Gotest}
				}
			case model.LintMessagesName:
				out, err = rc.taskrunner.FetchLintMessages(rc.linters)
				if err == nil {
					rc.Repository.SetLintMessages(out.(model.LintMessages))
				}
			}

			if err != nil {
				rc.Errors <- err
				rc.HasError = true
				return
			}

			rc.data <- out
		}(tp)

		lgr := rc.logger.WithField("type", tp)

		select {
		case err := <-rc.Errors:
			rc.Output[tp] = wrapError(err)
			lgr.Error(err)
			i++
		case out := <-rc.data:
			rc.Output[tp] = out
			i++
		case <-rc.Aborted:
			lgr.Warn("Shutting down (aborted)")
		case <-time.After(RoutineTimeout):
			rc.Output[tp] = wrapError(ErrRoutineTimeout)
			lgr.Error(ErrRoutineTimeout)
			i++
		}
	}

	// If every check has been ran
	if i == len(rc.types) {
		rc.StampEntry()

		// The entire dataset is ready
		rc.Done <- true

		go showcaser.ProcessRepository(rc.Repository)
	}
}

// StampEntry is called once the entire dataset is loaded.
func (rc *Checker) StampEntry() {
	// Add the metadata
	err := rc.Repository.SetMetadata()
	if err != nil {
		rc.Output[model.MetadataName] = wrapError(err)
	} else {
		d, _ := rc.Repository.GetMetadata()
		rc.Output[model.MetadataName] = d
	}

	// Add the score
	if rc.HasError {
		// If something went wrong during the processing
		// then we cannot calculate the rank
		rc.Output[model.ScoreName] = model.Score{Rank: ""}
	} else {
		err = rc.Repository.SetScore()
		if err != nil {
			rc.Output[model.ScoreName] = wrapError(err)
		} else {
			d, _ := rc.Repository.GetScore()
			rc.Output[model.ScoreName] = d
		}
	}

	// Add the timestamp
	err = rc.Repository.SetLastUpdate()
	if err != nil {
		rc.Output[model.LastUpdateName] = wrapError(err)
	} else {
		d, _ := rc.Repository.GetLastUpdate()
		rc.Output[model.LastUpdateName] = d
	}

	// Add the execution time
	err = rc.Repository.SetExecutionTime()
	if err != nil {
		rc.Output[model.ExecutionTimeName] = wrapError(err)
	} else {
		d, _ := rc.Repository.GetExecutionTime()
		rc.Output[model.ExecutionTimeName] = d
	}
}

// Abort declares the task as done and skips the processing.
func (rc *Checker) Abort() {
	rc.Aborted <- true
}

func wrapError(err error) interface{} {
	switch err := err.(type) {
	case processingError:
		return errorOutput{err.tp, err.message, err.output}
	}
	return errorOutput{Error: err.Error()}
}
