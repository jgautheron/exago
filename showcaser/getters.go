package showcaser

import "math/rand"

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
			Name:        repo.Name,
			Description: repo.Metadata.Description,
			Image:       repo.Metadata.Image,
			Rank:        repo.Score.Rank,
		})
	}
	return
}

func GetTopRankedRepositories() (repos []RepositoryInfo) {
	for i := 0; i <= data.itemCount; i++ {
		repo := data.topRanked[rand.Intn(len(data.topRanked))]
		repos = append(repos, RepositoryInfo{
			Name:        repo.Name,
			Description: repo.Metadata.Description,
			Image:       repo.Metadata.Image,
			Rank:        repo.Score.Rank,
		})
	}
	return
}

func GetPopularRepositories() (repos []RepositoryInfo) {
	for i := len(data.popular) - 1; i >= 0; i-- {
		repo := data.popular[i]
		repos = append(repos, RepositoryInfo{
			Name:        repo.Name,
			Description: repo.Metadata.Description,
			Image:       repo.Metadata.Image,
			Rank:        repo.Score.Rank,
		})
	}
	return
}
