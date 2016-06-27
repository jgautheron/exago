package mocks

import (
	"time"

	"github.com/exago/svc/repository"
	"github.com/exago/svc/repository/model"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	Name, Branch          string
	CodeStats             model.CodeStats
	Imports               model.Imports
	TestResults           model.TestResults
	LintMessages          model.LintMessages
	Metadata              model.Metadata
	Score                 model.Score
	StartTime, LastUpdate time.Time
	ExecutionTime         string

	repository.Repository
	mock.Mock
}

func NewRepositoryMock(repo, rank string) *RepositoryMock {
	return &RepositoryMock{
		Name: repo,
		Score: model.Score{
			Rank: rank,
		},
	}
}

func (r *RepositoryMock) GetName() string {
	return r.Name
}

func (r *RepositoryMock) GetRank() string {
	return r.Score.Rank
}

func (r *RepositoryMock) IsCached() bool {
	return true
}

func (r *RepositoryMock) IsLoaded() bool {
	return true
}

func (r *RepositoryMock) Load() (err error) {
	return nil
}

func (r *RepositoryMock) ClearCache() (err error) {
	return nil
}
