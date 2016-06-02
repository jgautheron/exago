package repository

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	// Lambda function time limit
	RoutineTimeout = time.Minute * 5
)

var (
	ErrRoutineTimeout = errors.New("The analysis timed out")
)

type Checker struct {
	logger         *log.Entry
	types, linters []string
	data           chan interface{}
	Repository     *Repository
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
		Done:       make(chan bool, 1),
		Output:     map[string]interface{}{},
	}
}

// Run launches concurrently every check and gathers the output.
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
			case "lintmessages":
				out, err = rc.Repository.GetLintMessages(rc.linters)
			}

			if err != nil {
				rc.data <- err
				return
			}

			rc.data <- out
		}(tp)

		lgr := rc.logger.WithField("type", tp)

		select {
		case out := <-rc.data:
			i++
			switch out.(type) {
			case error:
				err := out.(error)
				rc.Output[tp] = wrapError(err)
				lgr.Error(err)
			default:
				rc.Output[tp] = out
			}
		case <-time.After(RoutineTimeout):
			rc.Output[tp] = wrapError(ErrRoutineTimeout)
			lgr.Error(ErrRoutineTimeout)
		}
	}

	// If every check has been ran
	if i == len(rc.types) {
		rc.StampEntry()
	}
}

// StampEntry is called once the entire dataset is loaded.
func (rc *Checker) StampEntry() {
	// Add the score
	sc, err := rc.Repository.GetScore()
	if err != nil {
		rc.Output["score"] = wrapError(err)
	} else {
		rc.Output["score"] = sc
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

	// The entire dataset is ready
	rc.Done <- true
}

func wrapError(err error) interface{} {
	return struct {
		Error string `json:"error"`
	}{err.Error()}
}
