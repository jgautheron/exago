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
	Rank        string `json:"rank"`
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
			Rank:        "A",
		},
		{
			Name:        "github.com/inconshreveable/log15",
			Description: "Structured, composable logging for Go",
			Avatar:      "https://avatars1.githubusercontent.com/u/836692?v=3&s=460",
			Rank:        "C",
		},
		{
			Name:        "github.com/h2non/baloo",
			Description: "Expressive end-to-end HTTP API testing made easy in Go (golang)",
			Avatar:      "https://avatars1.githubusercontent.com/u/63402?v=3&s=460",
			Rank:        "B",
		},
		{
			Name:        "github.com/mholt/caddy",
			Description: "Fast, cross-platform HTTP/2 web server with automatic HTTPS",
			Avatar:      "https://avatars3.githubusercontent.com/u/1128849?v=3&s=460",
			Rank:        "C",
		},
		{
			Name:        "github.com/spf13/hugo",
			Description: "A Fast and Flexible Static Site Generator built with love in GoLang",
			Avatar:      "https://avatars0.githubusercontent.com/u/173412?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/hashicorp/consul",
			Description: "Consul is a tool for service discovery, monitoring and configuration.",
			Avatar:      "https://avatars0.githubusercontent.com/u/761456?v=3&s=200",
			Rank:        "F",
		},
	}
}

func GetTopRankedRepositories() []RepositoryInfo {
	return []RepositoryInfo{
		{
			Name:        "github.com/vova616/xxhash",
			Description: "xxhash is a pure go (golang) implementation of xxhash.",
			Avatar:      "https://avatars1.githubusercontent.com/u/1451914?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/clarkduvall/hyperloglog",
			Description: "HyperLogLog and HyperLogLog++ implementation in Go/Golang.",
			Avatar:      "https://avatars0.githubusercontent.com/u/1607111?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/spaolacci/murmur3",
			Description: "Native MurmurHash3 Go implementation",
			Avatar:      "https://avatars0.githubusercontent.com/u/3320120?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/jgautheron/codename-generator",
			Description: "A codename generator meant for naming software releases.",
			Avatar:      "https://avatars2.githubusercontent.com/u/683888?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/spf13/hugo",
			Description: "A Fast and Flexible Static Site Generator built with love in GoLang",
			Avatar:      "https://avatars0.githubusercontent.com/u/173412?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/Sirupsen/logrus",
			Description: "Structured, pluggable logging for Go.",
			Avatar:      "https://avatars2.githubusercontent.com/u/97400?v=3&s=460",
			Rank:        "A",
		},
	}
}

func GetPopularRepositories() []RepositoryInfo {
	return []RepositoryInfo{
		{
			Name:        "github.com/chreble/todo",
			Description: "A simple Go to-do list.",
			Avatar:      "https://avatars3.githubusercontent.com/u/62303?v=3&s=460",
			Rank:        "C",
		},
		{
			Name:        "github.com/reterVision/cuckoo",
			Description: "Cuckoo Filter implementation in Go",
			Avatar:      "https://avatars2.githubusercontent.com/u/1209592?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/dgryski/hokusai",
			Description: "Sketching streams in real-time",
			Avatar:      "https://avatars0.githubusercontent.com/u/970862?v=3&s=460",
			Rank:        "D",
		},
		{
			Name:        "github.com/rcrowley/go-metrics",
			Description: "Go port of Coda Hale's Metrics library",
			Avatar:      "https://avatars0.githubusercontent.com/u/11151?v=3&s=460",
			Rank:        "C",
		},
		{
			Name:        "github.com/dustin/go-probably",
			Description: "Probabilistic Data Structures for Go",
			Avatar:      "https://avatars1.githubusercontent.com/u/1779?v=3&s=460",
			Rank:        "D",
		},
		{
			Name:        "github.com/alicebob/miniredis",
			Description: "Pure Go Redis server for Go unittests",
			Avatar:      "https://avatars2.githubusercontent.com/u/621306?v=3&s=460",
			Rank:        "A",
		},
	}
}
