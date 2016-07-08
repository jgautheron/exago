package repository

import (
	"errors"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/exago/svc/github"
	"github.com/exago/svc/repository/model"
	"github.com/exago/svc/score"
)

func (r *Repository) SetImports(i model.Imports) {
	r.Imports = i
}

func (r *Repository) SetCodeStats(cs model.CodeStats) {
	r.CodeStats = cs
}

func (r *Repository) SetLintMessages(lm model.LintMessages) {
	r.LintMessages = lm
}

func (r *Repository) SetTestResults(tr model.TestResults) {
	r.TestResults = tr
}

func (r *Repository) SetStartTime(t time.Time) {
	r.StartTime = t
}

// SetExecutionTime sets the processing execution time.
// The value is then used to determine an ETA for refreshing data.
func (r *Repository) SetExecutionTime() (err error) {
	duration := time.Since(r.StartTime)
	r.ExecutionTime = (duration - (duration % time.Second)).String()
	return r.db.Put(r.cacheKey(model.ExecutionTimeName), []byte(r.ExecutionTime))
}

// SetLastUpdate sets the last update timestamp.
func (r *Repository) SetLastUpdate() (err error) {
	r.LastUpdate = time.Now()
	date := r.LastUpdate.Format(time.RFC3339)
	return r.db.Put(r.cacheKey(model.LastUpdateName), []byte(date))
}

// SetMetadata sets repository metadata such as description, stars...
func (r *Repository) SetMetadata() (err error) {
	reg, _ := regexp.Compile(`^github\.com/([\w\d\-]+)/([\w\d\-]+)`)
	m := reg.FindStringSubmatch(r.Name)
	if len(m) == 0 {
		return errors.New("Can only get metadata for GitHub repositories")
	}

	res, err := github.Get(m[1], m[2])
	if err != nil {
		return err
	}

	r.Metadata = model.Metadata{
		Image:       res["avatar_url"].(string),
		Description: res["description"].(string),
		Stars:       res["stargazers"].(int),
		LastPush:    res["last_push"].(time.Time),
	}

	return r.cacheData(model.MetadataName, r.Metadata)
}

// SetScore calculates the score based on the repository results.
func (r *Repository) SetScore() (err error) {
	val, res := score.Process(r.AsMap())
	r.Score.Value = val
	r.Score.Details = res
	r.Score.Rank = score.Rank(r.Score.Value)

	log.Infof(
		"[%s] Rank: %s, overall score: %.2f",
		r.GetName(),
		r.Score.Rank,
		r.Score.Value,
	)

	return nil
}
