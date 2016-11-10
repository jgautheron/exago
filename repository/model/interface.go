package model

import "time"

// Record represents a repository.
type Record interface {
	GetName() string
	GetBranch() string
	SetName(string)
	GetRank() string
	GetData() Data
	SetData(d Data)
	GetMetadata() Metadata
	SetMetadata(Metadata)
	GetLastUpdate() time.Time
	SetLastUpdate(t time.Time)
	GetExecutionTime() string
	SetExecutionTime(duration time.Duration)
	GetScore() Score
	GetProjectRunner() ProjectRunner
	SetProjectRunner(tr ProjectRunner)
	GetLintMessages() LintMessages
	SetLintMessages(LintMessages)
	SetError(tp string, err error)
	HasError() bool
	ApplyScore() (err error)
}

// RepositoryHost represents a VCS hosting service.
type RepositoryHost interface {
	GetFileContent(owner, repository, path string) (string, error)
	Get(owner, repository string) (map[string]interface{}, error)
}

// Load wraps repository-related methods.
type RepositoryLoader interface {
	Load(repository, branch string) (Record, error)
	IsValid(repository string) (map[string]interface{}, error)
	Save(repo Record) error
	ClearCache(repo, branch string) error
	IsCached(repo, branch string) bool
}

// Pool represents the routine pool.
type Pool interface {
	PushSync(repo string) (Record, error)
	PushAsync(repo string)
	WaitUntilEmpty()
	Stop()
}

// Promoter represents the promoting system mostly used to provide statistics.
type Promoter interface {
	AddRecent(repo Record)
	AddTopRanked(repo Record)
	AddPopular(repo Record)
	Process(repo Record)
	GetRecentRepositories() (repos []interface{})
	GetTopRankedRepositories() (repos []interface{})
	GetPopularRepositories() (repos []interface{})
}

// Database wraps the database calls.
type Database interface {
	DeleteAllMatchingPrefix(prefix []byte) error
	Delete(key []byte) error
	Put(key []byte, data []byte) error
	Get(key []byte) ([]byte, error)
}
