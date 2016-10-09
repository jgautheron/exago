package repository

import (
	"time"

	"github.com/hotolab/exago-svc/repository/model"
)

// GetName retrieves the full project name, including the provider domain name.
// Ex. github.com/hotolab/exago-svc
func (r *Repository) GetName() string {
	return r.Name
}

// GetRank retrieves the project's rank, ex "B+".
func (r *Repository) GetRank() string {
	return r.Data.Score.Rank
}

// GetMetadata retrieves repository metadata such as description, stars...
func (r *Repository) GetMetadata() model.Metadata {
	return r.Data.Metadata
}

// GetLastUpdate retrieves the timestamp when the project was last refreshed.
func (r *Repository) GetLastUpdate() time.Time {
	return r.Data.LastUpdate
}

// GetExecutionTime retrieves the last execution time.
// The value is used to determine an ETA for a project refresh.
func (r *Repository) GetExecutionTime() string {
	return r.Data.ExecutionTime
}

// GetScore retrieves the entire score details.
func (r *Repository) GetScore() model.Score {
	return r.Data.Score
}

// GetCodeStats retrieves the code statistics (LOC...).
func (r *Repository) GetCodeStats() model.CodeStats {
	return r.Data.CodeStats
}

// GetProjectRunner retrieves the test and checklist results.
func (r *Repository) GetProjectRunner() model.ProjectRunner {
	return r.Data.ProjectRunner
}

// GetLintMessages retrieves the linter warnings.
func (r *Repository) GetLintMessages() model.LintMessages {
	return r.Data.LintMessages
}

// GetData retrieves the repository data results.
func (r *Repository) GetData() model.Data {
	return r.Data
}

func (r *Repository) HasError() bool {
	return len(r.Data.Errors) > 0
}
