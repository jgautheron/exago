package processor

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	exago "github.com/hotolab/exago-svc"
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
	fns               = []string{"projectrunner", "lintmessages"}
)

type Processor struct {
	config exago.Config
}

type resultOutput struct {
	Fn       string
	Response job.Response
	err      error
}

func New(options ...exago.Option) *Processor {
	var p Processor
	for _, option := range options {
		option.Apply(&p.config)
	}
	return &p
}

func (p *Processor) ProcessRepository(value interface{}) interface{} {
	repo := value.(string)

	// Check first if the repository is valid (still exists, contains Go code...)
	data, err := p.config.RepositoryLoader.IsValid(repo)
	if err != nil {
		logger.WithField("repo", repo).Error(err)
		return err
	}

	startTime := time.Now()
	outCh := make(chan resultOutput, len(fns))
	wg := new(sync.WaitGroup)
	for _, fn := range fns {
		wg.Add(1)
		go func(fn, repo string) {
			defer wg.Done()
			out, err := job.CallLambdaFn(fn, repo, "")
			if err != nil {
				outCh <- resultOutput{
					Fn:  fn,
					err: err,
				}
				return
			}
			outCh <- resultOutput{fn, out, nil}
			logger.WithField("fn", fn).Debug("Received output")
		}(fn, repo)
	}
	wg.Wait()

	output := map[string]resultOutput{}
	for i := 0; i < len(fns); i++ {
		out := <-outCh
		output[out.Fn] = out
	}

	rp := p.importData(repo, output)
	rp.SetName(data["html_url"].(string))
	rp.SetMetadata(model.Metadata{
		Image:       data["avatar_url"].(string),
		Description: data["description"].(string),
		Stars:       data["stargazers"].(int),
		LastPush:    data["last_push"].(time.Time),
	})
	rp.SetExecutionTime(time.Since(startTime))
	rp.SetLastUpdate(time.Now())

	// Persist the dataset
	if err := p.config.RepositoryLoader.Save(rp); err != nil {
		logger.Errorf("Could not persist the dataset: %v", err)
	}

	return rp
}

func (p *Processor) importData(repo string, results map[string]resultOutput) model.Record {
	var err error
	rp := repository.New(repo, "")

	// Handle projectrunner
	var pr model.ProjectRunner
	if err = json.Unmarshal(*results[model.ProjectRunnerName].Response.Data, &pr); err != nil {
		rp.SetError(model.ProjectRunnerName, err)
	} else {
		rp.SetProjectRunner(pr)
	}

	// Handle lintmessages
	var lm model.LintMessages
	if err = json.Unmarshal(*results[model.LintMessagesName].Response.Data, &lm); err != nil {
		rp.SetError(model.LintMessagesName, err)
	} else {
		rp.SetLintMessages(lm)
	}

	// Add the score
	if err = rp.ApplyScore(); err != nil {
		rp.SetError(model.ScoreName, err)
	}

	return rp
}
