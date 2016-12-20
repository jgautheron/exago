package showcaser

import (
	"math/rand"
	"time"
)

type RepositoryInfo struct {
	Name        string `json:"name"`
	Branch      string `json:"branch"`
	GoVersion   string `json:"goversion"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Rank        string `json:"rank"`
}

func (s *Showcaser) GetRecentRepositories() (repos []interface{}) {
	for i := len(s.recent) - 1; i >= 0; i-- {
		repo := s.recent[i]
		meta := repo.GetMetadata()
		repos = append(repos, RepositoryInfo{
			Name:        repo.GetName(),
			Branch:      repo.GetBranch(),
			GoVersion:   repo.GetGoVersion(),
			Description: meta.Description,
			Image:       meta.Image,
			Rank:        repo.GetRank(),
		})
	}
	return
}

func (s *Showcaser) GetTopRankedRepositories() (repos []interface{}) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i, item := range rand.Perm(len(s.topRanked)) {
		if i == s.itemCount {
			break
		}
		repo := s.topRanked[item]
		meta := repo.GetMetadata()
		repos = append(repos, RepositoryInfo{
			Name:        repo.GetName(),
			Branch:      repo.GetBranch(),
			GoVersion:   repo.GetGoVersion(),
			Description: meta.Description,
			Image:       meta.Image,
			Rank:        repo.GetRank(),
		})
	}
	return
}

func (s *Showcaser) GetPopularRepositories() (repos []interface{}) {
	for _, repo := range s.popular {
		meta := repo.GetMetadata()
		repos = append(repos, RepositoryInfo{
			Name:        repo.GetName(),
			Branch:      repo.GetBranch(),
			GoVersion:   repo.GetGoVersion(),
			Description: meta.Description,
			Image:       meta.Image,
			Rank:        repo.GetRank(),
		})
	}
	return
}
