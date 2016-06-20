package repository

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/repository/lambda"
	"github.com/exago/svc/repository/model"
)

// GetDate retrieves the check timestamp.
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

// GetExecutionTime retrieves the last execution time.
// The value is used to determine an ETA for a project refresh.
func (r *Repository) GetExecutionTime() (string, error) {
	data, err := r.getCachedData("executiontime")
	if err != nil {
		return "", err
	}
	if data != nil {
		r.ExecutionTime = string(data)
		return r.ExecutionTime, nil
	}

	r.ExecutionTime = time.Since(r.StartTime).String()
	if err := leveldb.Save(r.cacheKey("executiontime"), []byte(r.ExecutionTime)); err != nil {
		return "", err
	}
	return r.ExecutionTime, nil
}

// GetScore retrieves the Exago score (A-F).
func (r *Repository) GetScore() (sc Score, err error) {
	data, err := r.getCachedData("score")
	if err != nil {
		return sc, err
	}
	if data != nil {
		if err := json.Unmarshal(data, &r.Score); err != nil {
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

// GetImports retrieves the third party imports.
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

		// Dedupe third party packages
		// One repository corresponds to one third party
		imports, filtered := []string{}, map[string]int{}
		reg, _ := regexp.Compile(`^([\w\d\.]+)/([\w\d\-]+)/([\w\d\-]+)`)
		for _, im := range r.Imports {
			m := reg.FindStringSubmatch(im)
			if len(m) > 0 {
				filtered[m[0]] = 1
			} else {
				filtered[im] = 1
			}
		}
		for im, _ := range filtered {
			imports = append(imports, im)
		}
		r.Imports = imports

		if err = r.cacheData(r.Imports.Name(), r.Imports); err != nil {
			return nil, err
		}
		return r.Imports, nil
	}
	if err := json.Unmarshal(data, &r.Imports); err != nil {
		return nil, err
	}
	return r.Imports, nil
}

// GetCodeStats retrieves the code statistics (LOC...).
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
	if err := json.Unmarshal(data, &r.CodeStats); err != nil {
		return nil, err
	}
	return r.CodeStats, nil
}

// GetTestResults retrieves the test and checklist results.
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
	if err := json.Unmarshal(data, &r.TestResults); err != nil {
		return tr, err
	}
	return r.TestResults, nil
}

// GetLintMessages retrieves the linter warnings emitted by gometalinter.
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
	if err := json.Unmarshal(data, &r.LintMessages); err != nil {
		return nil, err
	}
	return r.LintMessages, nil
}

// cacheKey formats the suffix as a standardised key.
func (r *Repository) cacheKey(suffix string) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s", r.Name, r.Branch, suffix))
}

// getCachedData attempts to load the data type from database.
func (r *Repository) getCachedData(suffix string) ([]byte, error) {
	return leveldb.FindForRepositoryCmd(r.cacheKey(suffix))
}

// cacheData persists the data type results in database.
func (r *Repository) cacheData(suffix string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return leveldb.Save(r.cacheKey(suffix), b)
}
