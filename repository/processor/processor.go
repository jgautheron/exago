package processor

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/pool/job"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/repository/model"
)

const (
	// Lambda function time limit
	RoutineTimeout = time.Second * 280
)

var (
	logger            = log.WithField("prefix", "processor")
	ErrRoutineTimeout = errors.New("The analysis timed out")
	fns               = []string{"codestats", "projectrunner", "lintmessages"}
)

type ResultOutput struct {
	Fn       string
	Response job.Response
	err      error
}

func ProcessRepository(value interface{}) interface{} {
	repo := value.(string)

	// Check first if the repository is valid (still exists, contains Go code...)
	data, err := repository.IsValid(repo)
	if err != nil {
		logger.WithField("repo", repo).Error(err)
		return err
	}

	startTime := time.Now()
	outCh := make(chan ResultOutput, len(fns))
	wg := new(sync.WaitGroup)
	for _, fn := range fns {
		wg.Add(1)
		go func(fn, repo string) {
			defer wg.Done()
			out, err := job.CallLambdaFn(fn, repo, "")
			if err != nil {
				outCh <- ResultOutput{
					Fn:  fn,
					err: err,
				}
				return
			}
			outCh <- ResultOutput{fn, out, nil}
			logger.Debugln(fn, out)
		}(fn, repo)
	}
	wg.Wait()

	output := map[string]ResultOutput{}
	for i := 0; i < len(fns); i++ {
		out := <-outCh
		output[out.Fn] = out
	}

	rp := importData(repo, output)
	rp.SetName(data["html_url"].(string))
	rp.SetExecutionTime(time.Since(startTime))
	rp.SetLastUpdate(time.Now())

	// Persist the dataset
	if err := rp.Save(); err != nil {
		logger.Errorf("Could not persist the dataset: %v", err)
	}

	return rp
}

func importData(repo string, results map[string]ResultOutput) *repository.Repository {
	var err error
	rp := repository.New(repo, "")

	// Handle codestats
	var cs model.CodeStats
	if err = json.Unmarshal(*results[model.CodeStatsName].Response.Data, &cs); err != nil {
		rp.SetError(model.CodeStatsName, err)
	} else {
		rp.SetCodeStats(cs)
	}

	// Handle projectrunner
	var pr model.ProjectRunner
	if err = json.Unmarshal(*results[model.ProjectRunnerName].Response.Data, &pr); err != nil {
		rp.SetError(model.ProjectRunnerName, err)
	} else {
		rp.SetProjectRunner(pr)
	}

	// Handle lintmessages
	var lm model.LintMessages
	logger.Warnln(string(*results[model.LintMessagesName].Response.Data))
	if err = json.Unmarshal(*results[model.LintMessagesName].Response.Data, &lm); err != nil {
		rp.SetError(model.LintMessagesName, err)
	} else {
		rp.SetLintMessages(lm)
	}

	// Add the metadata
	err = rp.SetMetadata()
	if err != nil {
		rp.SetError(model.MetadataName, err)
	}

	// Add the score
	err = rp.SetScore()
	if err != nil {
		rp.SetError(model.ScoreName, err)
	}

	return rp
}
