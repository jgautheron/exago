package processor

import (
	"errors"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/mocks"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/repository/model"
	. "github.com/stretchr/testify/mock"
)

var repo = "github.com/hotolab/foo"

func TestProcessed(t *testing.T) {
	rp := mocks.Record{Name: repo}
	rp.On("SetStartTime", Anything).Return(nil).Once()
	rp.On("Save").Return(nil).Once()

	tr := mocks.TaskRunner{}
	tr.On("FetchCodeStats").Return(model.CodeStats{}, nil).Once()
	rp.On("SetCodeStats", Anything).Return(nil).Once()

	tr.On("FetchProjectRunner").Return(model.ProjectRunner{}, nil).Once()
	rp.On("SetProjectRunner", Anything).Return(nil).Once()

	tr.On("FetchLintMessages", repository.DefaultLinters).Return(model.LintMessages{}, nil).Once()
	rp.On("SetLintMessages", Anything).Return(nil).Once()

	rp.On("SetMetadata").Return(nil).Once()
	rp.On("SetScore").Return(nil).Once()
	rp.On("SetLastUpdate").Return(nil).Once()
	rp.On("SetExecutionTime").Return(nil).Once()

	process := getMockProcessor(&rp, tr)
	process.Run()

	rp.AssertExpectations(t)
	tr.AssertExpectations(t)
}

func TestRunnerGotError(t *testing.T) {
	rp := mocks.Record{Name: repo}
	rp.On("SetStartTime", Anything).Return(nil).Once()
	rp.On("Save").Return(nil).Once()
	rp.On("SetError", Anything, Anything).Return(nil).Once()

	tr := mocks.TaskRunner{}
	tr.On("FetchCodeStats").Return(model.CodeStats{}, nil).Once()
	rp.On("SetCodeStats", Anything).Return(nil).Once()

	runner := model.ProjectRunner{}
	tr.On("FetchProjectRunner").Return(runner, errors.New("error")).Once()
	rp.On("SetProjectRunner", Anything).Return(nil).Once()

	tr.On("FetchLintMessages", repository.DefaultLinters).Return(model.LintMessages{}, nil).Once()
	rp.On("SetLintMessages", Anything).Return(nil).Once()

	rp.On("SetMetadata").Return(nil).Once()
	rp.On("SetScore").Return(nil).Once()
	rp.On("SetLastUpdate").Return(nil).Once()
	rp.On("SetExecutionTime").Return(nil).Once()

	process := getMockProcessor(&rp, tr)
	process.Run()
}

func TestProcessingAborted(t *testing.T) {
	rp := mocks.Record{Name: repo}
	rp.On("SetStartTime", Anything).Return(nil)

	tr := mocks.TaskRunner{}
	tr.On("FetchCodeStats").After(1*time.Second).Return(model.CodeStats{}, nil)

	process := getMockProcessor(&rp, tr)
	go func() {
		time.Sleep(500 * time.Millisecond)
		process.Abort()
	}()
	process.Run()

	rp.AssertExpectations(t)
	tr.AssertExpectations(t)
}

func getMockProcessor(rp *mocks.Record, tr mocks.TaskRunner) *Checker {
	return &Checker{
		logger:     log.WithField("repository", repo),
		types:      DefaultTypes,
		linters:    repository.DefaultLinters,
		taskrunner: tr,
		processed:  make(chan bool),
		Repository: rp,
		HasError:   false,
		Aborted:    make(chan bool, 1),
		Done:       make(chan bool, 1),
	}
}
