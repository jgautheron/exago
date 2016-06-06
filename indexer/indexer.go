// Package indexer processes repositories to determine occurrences, top-k, popularity.
package indexer

var data IndexedData

type IndexedData struct {
	Recent  []string
	Top     []string
	Popular []string
}

type RepositoryInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

func init() {
	data = IndexedData{}
}

func ProcessRepository(repo string) {

}

func GetRecentRepositories() []RepositoryInfo {
	return []RepositoryInfo{
		{
			Name:        "github.com/Sirupsen/logrus",
			Description: "Structured, pluggable logging for Go.",
			Avatar:      "https://avatars2.githubusercontent.com/u/97400?v=3&s=460",
		},
		{
			Name:        "github.com/inconshreveable/log15",
			Description: "Structured, composable logging for Go",
			Avatar:      "https://avatars1.githubusercontent.com/u/836692?v=3&s=460",
		},
		{
			Name:        "github.com/h2non/baloo",
			Description: "Expressive end-to-end HTTP API testing made easy in Go (golang)",
			Avatar:      "https://avatars1.githubusercontent.com/u/63402?v=3&s=460",
		},
		{
			Name:        "github.com/mholt/caddy",
			Description: "Fast, cross-platform HTTP/2 web server with automatic HTTPS",
			Avatar:      "https://avatars3.githubusercontent.com/u/1128849?v=3&s=460",
		},
		{
			Name:        "github.com/spf13/hugo",
			Description: "A Fast and Flexible Static Site Generator built with love in GoLang",
			Avatar:      "https://avatars0.githubusercontent.com/u/173412?v=3&s=460",
		},
		{
			Name:        "github.com/hashicorp/consul",
			Description: "Consul is a tool for service discovery, monitoring and configuration.",
			Avatar:      "https://avatars0.githubusercontent.com/u/761456?v=3&s=200",
		},
	}
}
