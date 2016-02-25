package github

import (
	"encoding/base64"

	"github.com/exago/svc/config"
	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	client *gh.Client

	// Short-hands
	Repositories *gh.RepositoriesService
)

func SetUp() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Values.GithubAccessToken},
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
