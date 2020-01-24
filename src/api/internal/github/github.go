package github

import (
	"context"
	"encoding/base64"
	"math/rand"
	"time"

	gh "github.com/google/go-github/github"
	"github.com/jgautheron/exago/src/api/internal/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type GitHub struct {
	accessTokenList  []string
	accessTokenIndex int
	client           *gh.Client
}

func NewWithConfig(ctx context.Context, cfg *config.GitHubConfig) (*GitHub, error) {
	accessTokenIndex := 0
	accessTokenList := cfg.GithubAccessTokens
	client := getClient(ctx, accessTokenList[accessTokenIndex])
	return &GitHub{accessTokenList, accessTokenIndex, client}, nil
}

func getClient(ctx context.Context, accessToken string) *gh.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)
	return client
}

// GetRateLimit helps keep track of the API rate limit.
func (g GitHub) GetRateLimit(ctx context.Context) (int, error) {
	limits, _, err := g.client.RateLimits(ctx)
	if err != nil {
		return 0, err
	}
	return limits.Core.Limit, nil
}

// DisplayRateLimit helps keeping track of GitHub's rate limits.
func (g GitHub) DisplayRateLimit(ctx context.Context) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	n := r.Intn(5)

	// Do not show every time
	if n != 0 && n%4 == 0 {
		limit, err := g.GetRateLimit(ctx)
		if err == nil {
			logrus.Debugf("Current rate limit: %d", limit)
		}
	}
}

// GetFileContent loads the file content of the given repository/filename.
func (g GitHub) GetFileContent(ctx context.Context, owner, repository, path string) (string, error) {
	file, _, _, err := g.repositories().GetContents(ctx, owner, repository, path, nil)
	if err != nil {
		return "", err
	}
	out := *file.Content

	g.DisplayRateLimit(ctx)

	if *file.Encoding == "base64" {
		b, err := base64.StdEncoding.DecodeString(out)
		if err != nil {
			return out, err
		}
		out = string(b)
	}
	return out, nil
}

// Get builds a trimmed down map of the few things we need to know about a repository.
func (g GitHub) Get(ctx context.Context, owner, repository string) (map[string]interface{}, error) {
	repo, _, err := g.repositories().Get(ctx, owner, repository)
	if err != nil {
		return nil, err
	}

	g.DisplayRateLimit(ctx)

	desc := ""
	if repo.Description != nil {
		desc = *repo.Description
	}

	languages, _, err := g.repositories().ListLanguages(ctx, owner, repository)
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
	return mp, nil
}

// repositories is a short-hand for the GitHub repos API.
// https://developer.github.com/v3/repos/
func (g GitHub) repositories() *gh.RepositoriesService {
	return g.client.Repositories
}

func (g *GitHub) RotateAccessToken() {
	nextIndex := g.accessTokenIndex + 1
	if nextIndex+1 > len(g.accessTokenList) {
		g.accessTokenIndex = nextIndex
	}
	g.client = getClient(context.Background(), g.accessTokenList[nextIndex])
}
