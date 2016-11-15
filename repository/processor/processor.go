package processor

import (
	"encoding/json"
	"errors"
	"strings"
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
	ErrRoutineTimeout = errors.New("The analysis timed out")
	ErrEmptyResponse  = errors.New("Empty response data")
	logger            = log.WithField("prefix", "processor")
	fns               = []string{"projectrunner", "lintmessages"}
)

type Processor struct {
	config exago.Config
}

type resultOutput struct {
	Repository, Branch, Fn string
	Response               job.Response
	err                    error
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
	branch := ""

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
		go func(fn, repo, branch string) {
			defer wg.Done()
			out, err := job.CallLambdaFn(fn, repo, branch)
			if err != nil {
				outCh <- resultOutput{
					Repository: repo,
					Branch:     branch,
					Fn:         fn,
					err:        err,
				}
				return
			}
			outCh <- resultOutput{repo, branch, fn, out, nil}
			logger.WithFields(log.Fields{
				"repository": repo,
				"branch":     branch,
				"fn":         fn,
			}).Debug("Received output")
		}(fn, repo, branch)
	}
	wg.Wait()

	output := map[string]resultOutput{}
	for i := 0; i < len(fns); i++ {
		out := <-outCh

		// Return directly the error if anything went wrong
		if out.err != nil {
			return out.err
		}

		output[out.Fn] = out
	}

	// Strip the protocol
	repositoryName := strings.Replace(data["html_url"].(string), "https://", "", 1)

	rp, err := p.importData(repo, branch, output)
	if err != nil {
		return err
	}

	rp.SetName(repositoryName)
	rp.SetMetadata(model.Metadata{
		Image:       data["avatar_url"].(string),
		Description: data["description"].(string),
		Stars:       data["stargazers"].(int),
		LastPush:    data["last_push"].(time.Time),
	})
	rp.SetExecutionTime(time.Since(startTime))
	rp.SetLastUpdate(time.Now())

	// Persist the dataset if everything went well
	if err := p.config.RepositoryLoader.Save(rp); err != nil {
		logger.Errorf("Could not persist the dataset: %v", err)
	}

	return rp
}

func (p *Processor) importData(repo, branch string, results map[string]resultOutput) (model.Record, error) {
	rp := repository.New(repo, "")

	// Handle projectrunner
	var pr model.ProjectRunner
	if err := extractData(results[model.ProjectRunnerName], &pr); err != nil {
		logError(repo, branch, fns[0], err)
		return nil, err
	}
	rp.SetProjectRunner(pr)

	// Handle lintmessages
	var lm model.LintMessages
	if err := extractData(results[model.LintMessagesName], &lm); err != nil {
		logError(repo, branch, fns[1], err)
		return nil, err
	}
	rp.SetLintMessages(lm)

	// Calculate the score
	if err := rp.ApplyScore(); err != nil {
		return nil, err
	}

	return rp, nil
}

func extractData(data resultOutput, out interface{}) error {
	if data.err != nil {
		return data.err
	}
	if data.Response.Data == nil {
		return ErrEmptyResponse
	}
	return json.Unmarshal(*data.Response.Data, &out)
}

func logError(repo, branch, fn string, err error) {
	logger.WithFields(log.Fields{
		"repository": repo,
		"branch":     branch,
		"fn":         fn,
	}).Error(err)
}
