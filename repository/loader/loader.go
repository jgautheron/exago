package loader

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	exago "github.com/hotolab/exago-svc"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/repository/model"
)

const (
	SizeLimit = 1500000 // bytes
)

var (
	ErrOnlyGitHub        = errors.New("Only GitHub is implemented right now")
	ErrInvalidRepository = errors.New("The repository name is invalid")
	ErrInvalidLanguage   = errors.New("The repository does not contain Go code")
	ErrTooLarge          = errors.New("The repository is too large")

	// Make sure it satisfies the interface.
	_ model.RepositoryLoader = (*Loader)(nil)
)

type Loader struct {
	config exago.Config
}

func New(options ...exago.Option) *Loader {
	var l Loader
	for _, option := range options {
		option.Apply(&l.config)
	}
	return &l
}

// Load retrieves the saved repository data from the database.
func (l Loader) Load(repo, branch, goversion string) (model.Record, error) {
	b, err := l.config.DB.Get(getCacheKey(repo, branch, goversion))
	if err != nil {
		return nil, err
	}

	var rp repository.Repository
	if err := json.Unmarshal(b, &rp); err != nil {
		return nil, err
	}

	return &rp, nil
}

func (l Loader) IsValid(repository string) (map[string]interface{}, error) {
	if !strings.HasPrefix(repository, "github.com") {
		return nil, ErrOnlyGitHub
	}
	re := regexp.MustCompile(`^github\.com/[\w\d\-\._]+/[\w\d\-\._]+`)
	if !re.MatchString(repository) {
		return nil, ErrInvalidRepository
	}

	sp := strings.Split(repository, "/")
	data, err := l.config.RepositoryHost.Get(sp[1], sp[2])
	if err != nil {
		return nil, err
	}

	languages := data["languages"].(map[string]int)
	size, exists := languages["Go"]
	if !exists {
		return nil, ErrInvalidLanguage
	}

	if size > SizeLimit {
		return nil, ErrTooLarge
	}

	return data, nil
}
