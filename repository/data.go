package repository

import (
	"fmt"
	"time"

	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/repository/lambda"
	"github.com/exago/svc/repository/model"
)

var (
	DefaultLinters = []string{
		"errcheck",
		"gofmt",
		"goimports",
		"golint",
		"deadcode",
		"dupl",
		"gocyclo",
		"ineffassign",
		"varcheck",
		"vet",
		"vetshadow",
	}
)

func (r *Repository) GetDate() (string, error) {
	data, err := r.getCachedData("date")
	if err != nil {
		return "", err
	}
	if data != nil {
		r.LastUpdate, _ = time.Parse(time.RFC3339, string(data))
		return string(data), nil
	}

	r.LastUpdate = time.Now()
	date := r.LastUpdate.Format(time.RFC3339)
	if err := leveldb.Save(r.cacheKey("date"), []byte(date)); err != nil {
		return "", err
	}
	return date, nil
}

func (r *Repository) GetScore() (sc Score, err error) {
	data, err := r.getCachedData("score")
	if err != nil {
		return sc, err
	}
	if data != nil {
		if err := msgpack.Unmarshal(data, &r.Score); err != nil {
			return sc, err
		}
		return r.Score, nil
	}
	sc = r.calcScore()
	if err = r.cacheData("score", sc); err != nil {
		return sc, err
	}
	return sc, nil
}

func (r *Repository) GetImports() (model.Imports, error) {
	data, err := r.getCachedData(r.Imports.Name())
	if err != nil {
		return nil, err
	}
	if data == nil {
		res, err := lambda.GetImports(r.Name)
		if err != nil {
			return nil, err
		}
		r.Imports = res.(model.Imports)
		if err = r.cacheData(r.Imports.Name(), r.Imports); err != nil {
			return nil, err
		}
		return r.Imports, nil
	}
	if err := msgpack.Unmarshal(data, &r.Imports); err != nil {
		return nil, err
	}
	return r.Imports, nil
}

func (r *Repository) GetCodeStats() (model.CodeStats, error) {
	var (
		data []byte
		err  error
	)

	data, err = r.getCachedData(r.CodeStats.Name())
	if err != nil {
		return nil, err
	}
	if data == nil {
		res, err := lambda.GetCodeStats(r.Name)
		if err != nil {
			return nil, err
		}
		r.CodeStats = res.(model.CodeStats)
		if err = r.cacheData(r.CodeStats.Name(), r.CodeStats); err != nil {
			return nil, err
		}
		return r.CodeStats, nil
	}
	if err := msgpack.Unmarshal(data, &r.CodeStats); err != nil {
		return nil, err
	}
	return r.CodeStats, nil
}

func (r *Repository) GetTestResults() (tr model.TestResults, err error) {
	data, err := r.getCachedData(r.TestResults.Name())
	if err != nil {
		return tr, err
	}
	if data == nil {
		res, err := lambda.GetTestResults(r.Name)
		if err != nil {
			return tr, err
		}
		r.TestResults = res.(model.TestResults)
		if err = r.cacheData(r.TestResults.Name(), r.TestResults); err != nil {
			return tr, err
		}
		return r.TestResults, nil
	}
	if err := msgpack.Unmarshal(data, &r.TestResults); err != nil {
		return tr, err
	}
	return r.TestResults, nil
}

func (r *Repository) GetLintMessages(linters []string) (model.LintMessages, error) {
	data, err := r.getCachedData(r.LintMessages.Name())
	if err != nil {
		return nil, err
	}
	if data == nil {
		res, err := lambda.GetLintMessages(r.Name, linters)
		if err != nil {
			return nil, err
		}
		r.LintMessages = res.(model.LintMessages)
		if err = r.cacheData(r.LintMessages.Name(), r.LintMessages); err != nil {
			return nil, err
		}
		return r.LintMessages, nil
	}
	if err := msgpack.Unmarshal(data, &r.LintMessages); err != nil {
		return nil, err
	}
	return r.LintMessages, nil
}

func (r *Repository) cacheKey(suffix string) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s", r.Name, r.Branch, suffix))
}

func (r *Repository) getCachedData(suffix string) ([]byte, error) {
	return leveldb.FindForRepositoryCmd(r.cacheKey(suffix))
}

func (r *Repository) cacheData(suffix string, data interface{}) error {
	b, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}
	return leveldb.Save(r.cacheKey(suffix), b)
}
