package processor

import (
	"errors"
	"strings"
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

func ProcessRepository(data interface{}) interface{} {
	outCh := make(chan ResultOutput, len(fns))

	repo := data.(string)
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
				logrus.Errorln(repo, fn, err)
				return
			}
			outCh <- ResultOutput{fn, out, nil}
		}(fn, repo)
	}
	wg.Wait()

	output := map[string]ResultOutput{}
	for i := 0; i < len(fns); i++ {
		out := <-outCh
		output[fns[i]] = out
	}

	return output
}