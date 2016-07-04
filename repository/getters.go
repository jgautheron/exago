package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/exago/svc/repository/model"
)

func (r *Repository) GetName() string {
	return r.Name
}

func (r *Repository) GetMetadataDescription() string {
	return r.Metadata.Description
}

func (r *Repository) GetMetadataImage() string {
	return r.Metadata.Image
}

func (r *Repository) GetMetadataStars() int {
	return r.Metadata.Stars
}

func (r *Repository) GetRank() string {
	return r.Score.Rank
}

// GetMetadata retrieves repository metadata such as description, stars...
func (r *Repository) GetMetadata() (d model.Metadata, err error) {
	data, err := r.getCachedData(model.MetadataName)
	if err != nil {
		return d, err
	}
	if err := json.Unmarshal(data, &r.Metadata); err != nil {
		return d, err
	}
	return r.Metadata, nil
}

// GetLastUpdate retrieves the last update timestamp.
func (r *Repository) GetLastUpdate() (string, error) {
	data, err := r.getCachedData(model.LastUpdateName)
	if err != nil {
		return "", err
	}
	r.LastUpdate, _ = time.Parse(time.RFC3339, string(data))
	return string(data), nil
}

// GetExecutionTime retrieves the last execution time.
// The value is used to determine an ETA for a project refresh.
func (r *Repository) GetExecutionTime() (string, error) {
	data, err := r.getCachedData(model.ExecutionTimeName)
	if err != nil {
		return "", err
	}
	r.ExecutionTime = string(data)
	return r.ExecutionTime, nil
}

// GetScore retrieves the Exago score (A-F).
func (r *Repository) GetScore() (sc model.Score, err error) {
	data, err := r.getCachedData(model.ScoreName)
	if err != nil {
		return sc, err
	}
	if err = json.Unmarshal(data, &r.Score); err != nil {
		return sc, err
	}
	return r.Score, nil
}

// GetImports retrieves the third party imports.
func (r *Repository) GetImports() (model.Imports, error) {
	data, err := r.getCachedData(model.ImportsName)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &r.Imports); err != nil {
		return nil, err
	}
	return r.Imports, nil
}

// GetCodeStats retrieves the code statistics (LOC...).
func (r *Repository) GetCodeStats() (model.CodeStats, error) {
	data, err := r.getCachedData(model.CodeStatsName)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &r.CodeStats); err != nil {
		return nil, err
	}
	return r.CodeStats, nil
}

// GetTestResults retrieves the test and checklist results.
func (r *Repository) GetTestResults() (tr model.TestResults, err error) {
	data, err := r.getCachedData(model.TestResultsName)
	if err != nil {
		return tr, err
	}
	if err := json.Unmarshal(data, &r.TestResults); err != nil {
		return tr, err
	}
	return r.TestResults, nil
}

// GetLintMessages retrieves the linter warnings emitted by gometalinter.
func (r *Repository) GetLintMessages(linters []string) (model.LintMessages, error) {
	data, err := r.getCachedData(model.LintMessagesName)
	if err != nil {
		return nil, err
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
	return r.db.FindForRepositoryCmd(r.cacheKey(suffix))
}

// cacheData persists the data type results in database.
func (r *Repository) cacheData(suffix string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.db.Put(r.cacheKey(suffix), b)
}
