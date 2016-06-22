package github

import (
	"encoding/base64"

	. "github.com/exago/svc/config"
	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	client *gh.Client

	// Short-hand
	Repositories *gh.RepositoriesService
)

func Init() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Config.GithubAccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = gh.NewClient(tc)

	Repositories = client.Repositories
}

func GetFileContent(owner, repository, path string) (string, error) {
	file, _, _, err := Repositories.GetContents(owner, repository, path, nil)
	if err != nil {
		return "", err
	}
	if *file.Encoding == "base64" {
		b, err := base64.StdEncoding.DecodeString(*file.Content)
		return string(b), err
	}
	return *file.Content, nil
}

func Get(owner, repository string) (map[string]interface{}, error) {
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

	// Export everything that could be useful
	return map[string]interface{}{
		"avatar_url":  *repo.Owner.AvatarURL,
		"description": desc,
		"language":    language,
		"stargazers":  *repo.StargazersCount,
		"last_push":   repo.PushedAt.Time,
	}, nil
}
