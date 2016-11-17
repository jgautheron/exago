package github

import (
	"encoding/base64"
	"math/rand"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	gh "github.com/google/go-github/github"
	"github.com/hashicorp/golang-lru"
	. "github.com/hotolab/exago-svc/config"
	"github.com/hotolab/exago-svc/repository/model"
	"golang.org/x/oauth2"
)

const (
	CacheSize = 20
)

var (
	logger = log.WithField("prefix", "github")

	// Make sure it satisfies the interface.
	_ model.RepositoryHost = (*GitHub)(nil)
)

type GitHub struct {
	client *gh.Client
	cache  *lru.ARCCache
	sync.RWMutex
}

func New() (GitHub, error) {
	// Authenticate with the GitHub API
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Config.GithubAccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := gh.NewClient(tc)

	// Initialise the ARC cache
	cache, err := lru.NewARC(CacheSize)
	if err != nil {
		return GitHub{}, err
	}

	return GitHub{
		client: client,
		cache:  cache,
	}, nil
}

// GetRateLimit helps keep track of the API rate limit.
func (g GitHub) GetRateLimit() (int, error) {
	limits, _, err := g.client.RateLimits()
	if err != nil {
		return 0, err
	}
	return limits.Core.Limit, nil
}

// DisplayRateLimit helps keeping track of GitHub's rate limits.
func (g GitHub) DisplayRateLimit() {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	n := r.Intn(5)

	// Do not show every time
	if n != 0 && n%4 == 0 {
		limit, err := g.GetRateLimit()
		if err == nil {
			logger.Debugf("Current rate limit: %d", limit)
		}
	}
}

// GetFileContent loads the file content of the given repository/filename.
func (g GitHub) GetFileContent(owner, repository, path string) (string, error) {
	if data, cached := g.getCached(owner, repository, path); cached {
		return data.(string), nil
	}
	file, _, _, err := g.repositories().GetContents(owner, repository, path, nil)
	if err != nil {
		return "", err
	}
	out := *file.Content

	g.DisplayRateLimit()

	if *file.Encoding == "base64" {
		b, err := base64.StdEncoding.DecodeString(out)
		if err != nil {
			return out, err
		}
		out = string(b)
	}
	g.saveCache(out, owner, repository, path)
	return out, nil
}

// Get builds a trimmed down map of the few things we need to know about a repository.
func (g GitHub) Get(owner, repository string) (map[string]interface{}, error) {
	if data, cached := g.getCached(owner, repository); cached {
		return data.(map[string]interface{}), nil
	}
	repo, _, err := g.repositories().Get(owner, repository)
	if err != nil {
		return nil, err
	}

	g.DisplayRateLimit()

	desc := ""
	if repo.Description != nil {
		desc = *repo.Description
	}

	languages, _, err := g.repositories().ListLanguages(owner, repository)
	if err != nil {
		return nil, err
	}

	// Export only what we need
	mp := map[string]interface{}{
		"html_url":    *repo.HTMLURL,
		"avatar_url":  *repo.Owner.AvatarURL,
		"description": desc,
		"languages":   languages,
		"stargazers":  *repo.StargazersCount,
		"last_push":   repo.PushedAt.Time,
	}
	g.saveCache(mp, owner, repository)
	return mp, nil
}

// repositories is a short-hand for the GitHub repos API.
// https://developer.github.com/v3/repos/
func (g GitHub) repositories() *gh.RepositoriesService {
	return g.client.Repositories
}
