package repository

import (
	"errors"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/exago/svc/repository/model"
	"github.com/exago/svc/score"
)

func (r *Repository) SetImports(i model.Imports) {
	r.Data.Imports = i
}

func (r *Repository) SetCodeStats(cs model.CodeStats) {
	r.Data.CodeStats = cs
}

func (r *Repository) SetLintMessages(lm model.LintMessages) {
	r.Data.LintMessages = lm
}

func (r *Repository) SetTestResults(tr model.TestResults) {
	r.Data.TestResults = tr
}

// SetStartTime stores the moment the processing started.
func (r *Repository) SetStartTime(t time.Time) {
	r.startTime = t
}

// SetExecutionTime sets the processing execution time.
// The value is then used to determine an ETA for refreshing data.
func (r *Repository) SetExecutionTime() {
	duration := time.Since(r.startTime)
	r.Data.ExecutionTime = (duration - (duration % time.Second)).String()
}

// SetLastUpdate sets the last update timestamp.
func (r *Repository) SetLastUpdate() {
	r.Data.LastUpdate = time.Now()
}

// SetMetadata sets repository metadata such as description, stars...
func (r *Repository) SetMetadata() (err error) {
	reg, _ := regexp.Compile(`^github\.com/([\w\d\-]+)/([\w\d\-]+)`)
	m := reg.FindStringSubmatch(r.Name)
	if len(m) == 0 {
		return errors.New("Can only get metadata for GitHub repositories")
	}

	res, err := r.RepositoryHost.Get(m[1], m[2])
	if err != nil {
		return err
	}

	r.Data.Metadata = model.Metadata{
		Image:       res["avatar_url"].(string),
		Description: res["description"].(string),
		Stars:       res["stargazers"].(int),
		LastPush:    res["last_push"].(time.Time),
	}

	return nil
}

// SetScore calculates the score based on the repository results.
func (r *Repository) SetScore() (err error) {
	val, res := score.Process(r.Data)
	r.Data.Score.Value = val
	r.Data.Score.Details = res
	r.Data.Score.Rank = score.Rank(r.Data.Score.Value)

	log.Infof(
		"[%s] Rank: %s, overall score: %.2f",
		r.GetName(),
		r.Data.Score.Rank,
		r.Data.Score.Value,
	)

	return nil
}

// SetError assigns a processing error to the given type (ex. TestResults).
// This helps keep track of what went wrong.
func (r *Repository) SetError(tp string, err error) {
	if r.Data.Errors == nil {
		r.Data.Errors = make(map[string]error)
	}
	r.Data.Errors[tp] = err
}
