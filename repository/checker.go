package repository

import (
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/indexer"
	"github.com/exago/svc/repository/model"
)

const (
	// Lambda function time limit
	RoutineTimeout = time.Minute * 5
)

var (
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

type Checker struct {
	logger         *log.Entry
	types, linters []string
	data           chan interface{}
	Repository     *Repository
	HasError       bool
	Errors         chan error
	Done           chan bool
	Output         map[string]interface{}
}

func NewChecker(repo string) *Checker {
	return &Checker{
		logger:     log.WithField("repository", repo),
		types:      DefaultTypes,
		linters:    DefaultLinters,
		data:       make(chan interface{}),
		Repository: New(repo, ""),
		HasError:   false,
		Errors:     make(chan error),
		Done:       make(chan bool, 1),
		Output:     map[string]interface{}{},
	}
}

// Run launches concurrently every check and merges the output.
func (rc *Checker) Run() {
	rc.Repository.StartTime = time.Now()

	i := 0
	for _, tp := range rc.types {
		go func(tp string) {
			var (
				out interface{}
				err error
			)

			switch tp {
			case "imports":
				out, err = rc.Repository.GetImports()
			case "codestats":
				out, err = rc.Repository.GetCodeStats()
			case "testresults":
				out, err = rc.Repository.GetTestResults()

				// Expose isolated errors
				switch ts := out.(model.TestResults); {
				case ts.Errors.Goget != "":
					err = processingError{"goget", ts.Errors.Goget, ts.RawOutput.Goget}
				case ts.Errors.Gotest != "":
					err = processingError{"gotest", ts.Errors.Gotest, ts.RawOutput.Gotest}
				}
			case "lintmessages":
				out, err = rc.Repository.GetLintMessages(rc.linters)
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

		go indexer.ProcessRepository(rc.Repository.Name)
	}
}

// StampEntry is called once the entire dataset is loaded.
func (rc *Checker) StampEntry() {
	if rc.HasError {
		rc.Output["score"] = Score{Rank: ""}
	} else {
		// Add the score
		sc, err := rc.Repository.GetScore()
		if err != nil {
			rc.Output["score"] = wrapError(err)
		} else {
			rc.Output["score"] = sc
		}
	}

	// Add the timestamp
	date, err := rc.Repository.GetDate()
	if err != nil {
		rc.Output["date"] = wrapError(err)
	} else {
		rc.Output["date"] = date
	}

	// Add the execution time
	et, err := rc.Repository.GetExecutionTime()
	if err != nil {
		rc.Output["execution_time"] = wrapError(err)
	} else {
		rc.Output["execution_time"] = et
	}
}

// Abort declares the task as done and skips the processing.
func (rc *Checker) Abort() {
	rc.Done <- true
}

func wrapError(err error) interface{} {
	switch err := err.(type) {
	case processingError:
		return errorOutput{err.tp, err.message, err.output}
	}
	return errorOutput{Error: err.Error()}
}
