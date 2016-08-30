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

func (s *Showcase) GetRecentRepositories() (repos []RepositoryInfo) {
	for i := len(s.recent) - 1; i >= 0; i-- {
		repo := s.recent[i]
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

func (s *Showcase) GetTopRankedRepositories() (repos []RepositoryInfo) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i, item := range rand.Perm(len(s.topRanked)) {
		if i == s.itemCount {
			break
		}
		repo := s.topRanked[item]
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

func (s *Showcase) GetPopularRepositories() (repos []RepositoryInfo) {
	for _, repo := range s.popular {
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
