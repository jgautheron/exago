package pool

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/pool/job"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/repository/processor"
	"github.com/hotolab/exago-svc/repository/processor"
	lm "github.com/hotolab/exago-svc/taskrunner/lambda"
	"github.com/jeffail/tunny"
)

const (
	SendTimeout = time.Second * 280
)

var (
	poolRunner *PoolRunner
	once       sync.Once
	fns        = []string{"codestats", "projectrunner", "lintmessages"}
)

type PoolRunner struct {
	pool        *tunny.WorkPool
	processorFn func(data interface{}) interface{}
}

type ResultOutput struct {
	Fn       string
	Response lm.Response
	err      error
}

// GetInstance returns the queue instance.
// The queue will be instantiated if it wasn't yet.
func GetInstance() *PoolRunner {
	once.Do(func() {
		numCPUs := 4
		pool, _ := tunny.CreatePool(numCPUs, ProcessRepository).Open()
		poolRunner = &PoolRunner{
			pool:        pool,
			processorFn: processor.ProcessRepository,
		}
	})
	return poolRunner
}

func (pr *PoolRunner) PushSync(repo string) interface{} {
	value, _ := pr.pool.SendWork(repo)
	return value
}

func (pr *PoolRunner) PushAsync(repo string) {
	pr.pool.SendWorkAsync(repo, nil)
}

func (pr *PoolRunner) Stop() {
	pr.pool.Close()
}

func TestCustomWorkers() {
	job.Init()

	pr := GetInstance()

	repos := []string{
		"github.com/Arachnid/evmdis",
		"github.com/mailhog/http",
		"github.com/stianeikeland/go-rpio",
		"github.com/codahale/safecookie",
		"github.com/alecthomas/jsonschema",
		"github.com/Synthace/go-glpk",
		"github.com/bezrukovspb/mux",
		"github.com/levicook/smaug",
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(repos))
	for i := 0; i < len(repos); i++ {
		go func(i int) {
			startTime := time.Now()

			repo := repos[i]
			value := pr.PushSync(repo)

			duration := time.Since(startTime)
			logrus.Infoln("PROCESSED", repo)

			rp := loadRepo(repo, value.(map[string]ResultOutput))
			rp.SetExecutionTime(duration)
			rp.SetLastUpdate(time.Now())

			// Persist the dataset
			// if err := rp.Save(); err != nil {
			// 	rc.logger.Errorf("Could not persist the dataset: %v", err)
			// }

			wg.Done()
		}(i)
	}

	wg.Wait()
}

func loadRepo(repo string, results map[string]ResultOutput) *repository.Repository {
	rp := repository.New(repo, "")

	// Handle codestats
	var cs model.CodeStats
	if err := json.Unmarshal(*results[model.CodeStatsName].Response.Data, &cs); err != nil {
		rp.SetError(model.CodeStatsName, err)
	} else {
		rp.SetCodeStats(cs)
	}

	// Handle projectrunner
	var pr model.ProjectRunner
	if err := json.Unmarshal(*results[model.ProjectRunnerName].Response.Data, &pr); err != nil {
		rp.SetError(model.ProjectRunnerName, err)
	} else {
		rp.SetProjectRunner(pr)
	}

	// Handle lintmessages
	var lm model.LintMessages
	if err := json.Unmarshal(*results[model.LintMessagesName].Response.Data, &lm); err != nil {
		rp.SetError(model.LintMessagesName, err)
	} else {
		rp.SetLintMessages(lm)
	}

	// Add the metadata
	err := rp.SetMetadata()
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
