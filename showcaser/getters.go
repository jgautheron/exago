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
		meta := repo.GetMetadata()
		repos = append(repos, RepositoryInfo{
			Name:        repo.GetName(),
			Description: meta.Description,
			Image:       meta.Image,
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
		meta := repo.GetMetadata()
		repos = append(repos, RepositoryInfo{
			Name:        repo.GetName(),
			Description: meta.Description,
			Image:       meta.Image,
			Rank:        repo.GetRank(),
		})
	}
	return
}

func GetPopularRepositories() (repos []RepositoryInfo) {
	for _, repo := range data.popular {
		meta := repo.GetMetadata()
		repos = append(repos, RepositoryInfo{
			Name:        repo.GetName(),
			Description: meta.Description,
			Image:       meta.Image,
			Rank:        repo.GetRank(),
		})
	}
	return
}
