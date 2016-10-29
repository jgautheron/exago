package github

import (
	"encoding/base64"
	"sync"

	log "github.com/Sirupsen/logrus"
	gh "github.com/google/go-github/github"
	"github.com/hashicorp/golang-lru"
	. "github.com/hotolab/exago-svc/config"
	"golang.org/x/oauth2"
)

const (
	CacheSize = 20
)

var (
	gi   GitHub
	once sync.Once

	// Make sure it satisfies the interface.
	_ RepositoryHost = (*GitHub)(nil)
)

type GitHub struct {
	client *gh.Client
	cache  *lru.ARCCache
	sync.RWMutex
}

func GetInstance() GitHub {
	once.Do(func() {
		// Authenticate with the GitHub API
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: Config.GithubAccessToken},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		client := gh.NewClient(tc)

		// Initialise the ARC cache
		cache, err := lru.NewARC(CacheSize)
		if err != nil {
			log.Fatal("Could not create the ARC cache")
		}

		gi = GitHub{
			client: client,
			cache:  cache,
		}
	})
	return gi
}

func (g GitHub) GetFileContent(owner, repository, path string) (string, error) {
	if data, cached := g.getCached(owner, repository, path); cached {
		return data.(string), nil
	}
	file, _, _, err := g.repositories().GetContents(owner, repository, path, nil)
	if err != nil {
		return "", err
	}
	out := *file.Content

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

func (g GitHub) Get(owner, repository string) (map[string]interface{}, error) {
	if data, cached := g.getCached(owner, repository); cached {
		return data.(map[string]interface{}), nil
	}
	repo, _, err := g.repositories().Get(owner, repository)
	if err != nil {
		return nil, err
	}

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

type RepositoryHost interface {
	GetFileContent(owner, repository, path string) (string, error)
	Get(owner, repository string) (map[string]interface{}, error)
}
