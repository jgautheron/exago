package model

import "time"

// Record represents a repository.
type Record interface {
	GetName() string
	SetName(string)
	GetRank() string
	GetData() Data
	GetMetadata() Metadata
	SetMetadata(Metadata)
	GetLastUpdate() time.Time
	SetLastUpdate(t time.Time)
	GetExecutionTime() string
	SetExecutionTime(duration time.Duration)
	GetScore() Score
	GetCodeStats() CodeStats
	SetCodeStats(CodeStats)
	GetProjectRunner() ProjectRunner
	SetProjectRunner(tr ProjectRunner)
	GetLintMessages() LintMessages
	SetLintMessages(LintMessages)
	SetStartTime(t time.Time)
	SetError(tp string, err error)
	HasError() bool
	ApplyScore() (err error)
	IsCached() bool
	IsLoaded() bool
	Load() (err error)
	ClearCache() (err error)
	Save() error
}

// RepositoryHost represents a VCS hosting service.
type RepositoryHost interface {
	GetFileContent(owner, repository, path string) (string, error)
	Get(owner, repository string) (map[string]interface{}, error)
}

// Load wraps repository-related methods.
type RepositoryLoader interface {
	Load(repository, branch string) Record
	IsValid(repository string) (map[string]interface{}, error)
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
