package processor

import (
	"errors"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/taskrunner"
	"github.com/hotolab/exago-svc/taskrunner/lambda"
)

const (
	// Lambda function time limit
	RoutineTimeout = time.Second * 280
)

var (
	logger            = log.WithField("prefix", "processor")
	ErrRoutineTimeout = errors.New("The analysis timed out")
)

func ProcessRepository(repo, branch string, tr taskrunner.TaskRunner) (interface{}, error) {
	// data, err := repository.IsValid(repo)
	// if err != nil {
	// 	return nil, err
	// }

	// Real repository name (with caps if any)
	// rp := strings.Replace(data["html_url"].(string), "https://", "", 1)

	rp := repo
	checker := NewChecker(rp, branch, tr)
	checker.Run()

	var out model.Data
	select {
	case <-checker.Done:
		out = checker.Repository.GetData()
		out.Name = rp
		out.Branch = branch
	case <-checker.Aborted:
		out = checker.Repository.GetData()
		out.Name = rp
		out.Branch = branch
	default:
	}
	return out, nil
}

type Checker struct {
	logger     *log.Entry
	taskrunner taskrunner.TaskRunner
	Repository repository.Repository
	Aborted    chan bool
	Done       chan bool
}

func NewChecker(repo, branch string, tr taskrunner.TaskRunner) *Checker {
	return &Checker{
		logger:     logger.WithField("repository", repo),
		taskrunner: tr,
		Repository: repository.Repository{Name: repo},
		Aborted:    make(chan bool, 1),
		Done:       make(chan bool, 1),
	}
}

// Run launches concurrently every check and merges the output.
func (rc *Checker) Run() {
	// rc.Repository.SetStartTime(time.Now())

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(rc *Checker) {
		rc.logger.Warnln("start FetchCodeStats")
		defer wg.Done()
		tr := lambda.Runner{Repository: rc.Repository.Name}
		out, _ := tr.FetchCodeStats()
		rc.Repository.Data.CodeStats = out
		// time.Sleep(10 * time.Second)
		rc.logger.Warnln("end FetchCodeStats")
	}(rc)

	wg.Wait()

	wg.Add(1)
	go func(rc *Checker) {
		rc.logger.Warnln("start FetchProjectRunner")
		defer wg.Done()
		tr := lambda.Runner{Repository: rc.Repository.Name}
		out, _ := tr.FetchProjectRunner()
		rc.Repository.Data.ProjectRunner = out
		// time.Sleep(10 * time.Second)
		rc.logger.Warnln("end FetchProjectRunner")
	}(rc)
	wg.Wait()

	wg.Add(1)
	go func(rc *Checker) {
		rc.logger.Warnln("start FetchLintMessages")
		defer wg.Done()
		tr := lambda.Runner{Repository: rc.Repository.Name}
		out, _ := tr.FetchLintMessages()
		rc.Repository.Data.LintMessages = out
		// time.Sleep(10 * time.Second)
		rc.logger.Warnln("end FetchLintMessages")
	}(rc)

	wg.Wait()
	// rc.StampEntry()
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
