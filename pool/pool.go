package pool

import (
	"time"

	"github.com/jeffail/tunny"
	exago "github.com/jgautheron/exago"
	. "github.com/jgautheron/exago/config"
	"github.com/jgautheron/exago/repository/model"
)

const (
	SendTimeout = time.Second * 280
)

var (
	// Make sure it satisfies the interface.
	_ model.Pool = (*PoolRunner)(nil)
)

type PoolRunner struct {
	pool   *tunny.WorkPool
	config exago.Config
}

func New(options ...exago.Option) (model.Pool, error) {
	var p PoolRunner
	for _, option := range options {
		option.Apply(&p.config)
	}
	pool, err := tunny.CreatePool(Config.PoolSize, p.config.RepositoryProcessor).Open()
	if err != nil {
		return nil, err
	}
	p.pool = pool
	return &p, nil
}

func (pr *PoolRunner) PushSync(repo string) (model.Record, error) {
	value, _ := pr.pool.SendWork(repo)
	switch value.(type) {
	case error:
		return nil, value.(error)
	default:
		return value.(model.Record), nil
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
