package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/exago/svc/leveldb"
)

var (
	errMissingData = errors.New("Not enough data to calculate the rank")
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

	passed := 0
	for k, v := range data {
		switch k.Type {
		case "codestats":
			if err := json.Unmarshal(v, &r.CodeStats); err != nil {
				return err
			}
			passed++
		case "imports":
			if err := json.Unmarshal(v, &r.Imports); err != nil {
				return err
			}
			passed++
		case "testrunner":
			if err := json.Unmarshal(v, &r.TestResults); err != nil {
				return err
			}
			passed++
		}
	}

	// All three are required to determine the rank
	if passed != 3 {
		return errMissingData
	}

	r.calcScore()
	return nil
}

func (r *Repository) Rank() Rank {
	return r.Score.Rank
}
