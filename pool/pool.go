package pool

import (
	"sync"
	"time"

	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/repository/processor"
	"github.com/jeffail/tunny"
)

const (
	SendTimeout = time.Second * 280
)

var (
	poolRunner *PoolRunner
	once       sync.Once
)

type PoolRunner struct {
	pool *tunny.WorkPool
}

// GetInstance returns the pool instance.
func GetInstance() *PoolRunner {
	once.Do(func() {
		numCPUs := 4
		pool, _ := tunny.CreatePool(numCPUs, processor.ProcessRepository).Open()
		poolRunner = &PoolRunner{
			pool: pool,
		}
	})
	return poolRunner
}

func (pr *PoolRunner) PushSync(repo string) (*repository.Repository, error) {
	value, _ := pr.pool.SendWork(repo)
	switch value.(type) {
	case error:
		return nil, value.(error)
	default:
		return value.(*repository.Repository), nil
	}
	return nil, nil
}

func (pr *PoolRunner) PushAsync(repo string) {
	pr.pool.SendWorkAsync(repo, nil)
}

func (pr *PoolRunner) WaitUntilEmpty() {
	for {
		time.Sleep(1 * time.Second)
		if pr.pool.NumPendingAsyncJobs() == 0 {
			return
		}
	}
}

func (pr *PoolRunner) Stop() {
	pr.pool.Close()
}
