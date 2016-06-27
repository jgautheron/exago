package github

import (
	"encoding/base64"

	. "github.com/exago/svc/config"
	gh "github.com/google/go-github/github"
	"github.com/hashicorp/golang-lru"
	"golang.org/x/oauth2"
)

const (
	CacheSize = 20
)

var (
	client *gh.Client
	cache  *lru.ARCCache

	// Repositories is a short-hand for the GitHub repos API.
	// https://developer.github.com/v3/repos/
	Repositories *gh.RepositoriesService
)

func Init() (err error) {
	// Authenticate with the GitHub API
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Config.GithubAccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = gh.NewClient(tc)
	Repositories = client.Repositories

	// Initialise the ARC cache
	cache, err = lru.NewARC(CacheSize)
	return err
}

func GetFileContent(owner, repository, path string) (string, error) {
	if data, cached := getCached(owner, repository, path); cached {
		return data.(string), nil
	}
	file, _, _, err := Repositories.GetContents(owner, repository, path, nil)
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
	saveCache(out, owner, repository, path)
	return out, nil
}

func Get(owner, repository string) (map[string]interface{}, error) {
	if data, cached := getCached(owner, repository); cached {
		return data.(map[string]interface{}), nil
	}
	repo, _, err := Repositories.Get(owner, repository)
	if err != nil {
		return nil, err
	}

	language := ""
	if repo.Language != nil {
		language = *repo.Language
	}

	desc := ""
	if repo.Description != nil {
		desc = *repo.Description
	}

	// Export only what we need
	mp := map[string]interface{}{
		"avatar_url":  *repo.Owner.AvatarURL,
		"description": desc,
		"language":    language,
		"stargazers":  *repo.StargazersCount,
		"last_push":   repo.PushedAt.Time,
	}
	saveCache(mp, owner, repository)
	return mp, nil
}
