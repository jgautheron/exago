package repository

import (
	"errors"
	"regexp"
	"strings"

	"github.com/hotolab/exago-svc/github"
)

const (
	SizeLimit = 1000000 // bytes
)

var (
	ErrOnlyGitHub        = errors.New("Only GitHub is implemented right now")
	ErrInvalidRepository = errors.New("The repository name is invalid")
	ErrInvalidLanguage   = errors.New("The repository does not contain Go code")
	ErrTooLarge          = errors.New("The repository is too large")
)

func IsValid(repository string) (map[string]interface{}, error) {
	if !strings.HasPrefix(repository, "github.com") {
		return nil, ErrOnlyGitHub
	}
	re := regexp.MustCompile(`^github\.com/[\w\d\-\._]+/[\w\d\-\._]+`)
	if !re.MatchString(repository) {
		return nil, ErrInvalidRepository
	}

	sp := strings.Split(repository, "/")
	data, err := github.GetInstance().Get(sp[1], sp[2])
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
