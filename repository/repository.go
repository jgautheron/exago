package repository

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/exago/svc/leveldb"
)

type Repository struct {
	name, branch string

	CodeStats    CodeStats
	Imports      Imports
	TestResults  TestResults
	LintMessages LintMessages

	Score Score
}

func New(repo, branch string) *Repository {
	return &Repository{
		name:   repo,
		branch: branch,
	}
}

func (r *Repository) LoadFromDB() error {
	prefix := fmt.Sprintf("%s-%s", r.name, r.branch)
	data, err := leveldb.FindAllForRepository([]byte(prefix))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("No data found in database for this repository")
	}

	for k, v := range data {
		logrus.Infoln(k.Type, v)
	}

	return nil
}

func (r *Repository) Rank() (rnk Rank, err error) {
	err = r.LoadFromDB()
	return rnk, err
}
