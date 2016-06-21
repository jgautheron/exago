package indexer

type RepositoryInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"avatar"`
	Rank        string `json:"rank"`
}

func GetRecentRepositories() []RepositoryInfo {
	var repos []RepositoryInfo
	for _, repo := range data.recent {
		repos = append(repos, RepositoryInfo{
			Name:        repo.Name,
			Description: repo.Metadata.Description,
			Image:       repo.Metadata.Image,
			Rank:        repo.Score.Rank,
		})
	}
	return repos
}

func GetTopRankedRepositories() []RepositoryInfo {
	return []RepositoryInfo{
		{
			Name:        "github.com/vova616/xxhash",
			Description: "xxhash is a pure go (golang) implementation of xxhash.",
			Image:       "https://avatars1.githubusercontent.com/u/1451914?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/clarkduvall/hyperloglog",
			Description: "HyperLogLog and HyperLogLog++ implementation in Go/Golang.",
			Image:       "https://avatars0.githubusercontent.com/u/1607111?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/spaolacci/murmur3",
			Description: "Native MurmurHash3 Go implementation",
			Image:       "https://avatars0.githubusercontent.com/u/3320120?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/jgautheron/codename-generator",
			Description: "A codename generator meant for naming software releases.",
			Image:       "https://avatars2.githubusercontent.com/u/683888?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/spf13/hugo",
			Description: "A Fast and Flexible Static Site Generator built with love in GoLang",
			Image:       "https://avatars0.githubusercontent.com/u/173412?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/Sirupsen/logrus",
			Description: "Structured, pluggable logging for Go.",
			Image:       "https://avatars2.githubusercontent.com/u/97400?v=3&s=460",
			Rank:        "A",
		},
	}
}

func GetPopularRepositories() []RepositoryInfo {
	return []RepositoryInfo{
		{
			Name:        "github.com/chreble/todo",
			Description: "A simple Go to-do list.",
			Image:       "https://avatars3.githubusercontent.com/u/62303?v=3&s=460",
			Rank:        "C",
		},
		{
			Name:        "github.com/reterVision/cuckoo",
			Description: "Cuckoo Filter implementation in Go",
			Image:       "https://avatars2.githubusercontent.com/u/1209592?v=3&s=460",
			Rank:        "A",
		},
		{
			Name:        "github.com/dgryski/hokusai",
			Description: "Sketching streams in real-time",
			Image:       "https://avatars0.githubusercontent.com/u/970862?v=3&s=460",
			Rank:        "D",
		},
		{
			Name:        "github.com/rcrowley/go-metrics",
			Description: "Go port of Coda Hale's Metrics library",
			Image:       "https://avatars0.githubusercontent.com/u/11151?v=3&s=460",
			Rank:        "C",
		},
		{
			Name:        "github.com/dustin/go-probably",
			Description: "Probabilistic Data Structures for Go",
			Image:       "https://avatars1.githubusercontent.com/u/1779?v=3&s=460",
			Rank:        "D",
		},
		{
			Name:        "github.com/alicebob/miniredis",
			Description: "Pure Go Redis server for Go unittests",
			Image:       "https://avatars2.githubusercontent.com/u/621306?v=3&s=460",
			Rank:        "A",
		},
	}
}
