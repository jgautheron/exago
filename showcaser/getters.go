package showcaser

import (
	"math/rand"
	"time"
)

type RepositoryInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Rank        string `json:"rank"`
}

func GetRecentRepositories() (repos []RepositoryInfo) {
	for i := len(data.recent) - 1; i >= 0; i-- {
		repo := data.recent[i]
		repos = append(repos, RepositoryInfo{
			Name:        repo.GetName(),
			Description: repo.GetMetadataDescription(),
			Image:       repo.GetMetadataImage(),
			Rank:        repo.GetRank(),
		})
	}
	return
}

func GetTopRankedRepositories() (repos []RepositoryInfo) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i, item := range rand.Perm(len(data.topRanked)) {
		if i == data.itemCount {
			break
		}
		repo := data.topRanked[item]
		repos = append(repos, RepositoryInfo{
			Name:        repo.GetName(),
			Description: repo.GetMetadataDescription(),
			Image:       repo.GetMetadataImage(),
			Rank:        repo.GetRank(),
		})
	}
	return
}

func GetPopularRepositories() (repos []RepositoryInfo) {
	for _, repo := range data.popular {
		repos = append(repos, RepositoryInfo{
			Name:        repo.GetName(),
			Description: repo.GetMetadataDescription(),
			Image:       repo.GetMetadataImage(),
			Rank:        repo.GetRank(),
		})
	}
	return
}
