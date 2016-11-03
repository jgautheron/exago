package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/hotolab/exago-svc/badge"
	"github.com/hotolab/exago-svc/godoc"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) repositoryHandler(w http.ResponseWriter, r *http.Request) {
	repo := context.Get(r, "repository").(string)
	branch := ""
	if s.config.RepositoryLoader.IsCached(repo, branch) {
		rp, err := s.config.RepositoryLoader.Load(repo, branch)
		if err != nil {
			send(w, r, rp, err)
		}
		if !rp.HasError() {
			go s.config.Showcaser.Process(rp)
		}
		send(w, r, rp.GetData(), nil)
		return
	}

	rp, err := s.config.Pool.PushSync(repo)
	if err == nil && !rp.HasError() {
		go s.config.Showcaser.Process(rp)
	}
	send(w, r, rp.GetData(), err)
}

func (s *Server) refreshHandler(w http.ResponseWriter, r *http.Request) {
	repo := context.Get(r, "repository").(string)
	branch := ""
	s.config.RepositoryLoader.ClearCache(repo, branch)
	s.repositoryHandler(w, r)
}

func (s *Server) badgeHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	lgr := context.Get(r, "lgr").(*log.Entry)

	repo := context.Get(r, "repository").(string)
	branch := ""

	if !s.config.RepositoryLoader.IsCached(repo, branch) {
		badge.WriteError(w, "")
		return
	}

	rp, err := s.config.RepositoryLoader.Load(repo, branch)
	if err != nil {
		lgr.Error(err)
		badge.WriteError(w, "")
		return
	}
	switch tp := ps.ByName("type"); tp {
	case "rank":
		rank, score := rp.GetRank(), rp.GetScore().Value

		badge.Write(w, "", string(rank), score)
	case "cov":
		title := "coverage"
		cov, score := rp.GetProjectRunner().GetMeanCodeCov(), rp.GetScore().Value
		badge.Write(w, title, fmt.Sprintf("%.2f%%", cov), score)
	case "duration":
		title := "tests duration"
		avg, score := rp.GetProjectRunner().GetMeanTestDuration(), rp.GetScore().Value
		badge.Write(w, title, fmt.Sprintf("%.2fs", avg), score)
	case "tests":
		title := "tests"
		tests, score := rp.GetCodeStats()["Test"], rp.GetScore().Value
		badge.Write(w, title, fmt.Sprintf("%d", tests), score)
	case "thirdparties":
		title := "3rd parties"
		thirdParties, score := len(rp.GetProjectRunner().Thirdparties.Data), rp.GetScore().Value
		badge.Write(w, title, fmt.Sprintf("%d", thirdParties), score)
	case "loc":
		title := "LOC"
		thirdParties, score := rp.GetCodeStats()["LOC"], rp.GetScore().Value
		badge.Write(w, title, fmt.Sprintf("%d", thirdParties), score)
	}
}

func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request) {
	owner := context.Get(r, "owner").(string)
	project := context.Get(r, "project").(string)
	path := context.Get(r, "path").(string)
	content, err := s.config.RepositoryHost.GetFileContent(owner, project, path)
	send(w, r, content, err)
}

func (s *Server) cachedHandler(w http.ResponseWriter, r *http.Request) {
	repo := context.Get(r, "repository").(string)
	branch := ""
	if s.config.RepositoryLoader.IsCached(repo, branch) {
		rp, err := s.config.RepositoryLoader.Load(repo, branch)
		send(w, r, rp.GetData(), err)
		return
	}
	send(w, r, false, nil)
}

func (s *Server) recentHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.config.Showcaser.GetRecentRepositories()
	out := map[string]interface{}{
		"type":         "recent",
		"repositories": repos,
	}
	send(w, r, out, nil)
}

func (s *Server) topHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.config.Showcaser.GetTopRankedRepositories()
	out := map[string]interface{}{
		"type":         "top",
		"repositories": repos,
	}
	send(w, r, out, nil)
}

func (s *Server) popularHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.config.Showcaser.GetPopularRepositories()
	out := map[string]interface{}{
		"type":         "popular",
		"repositories": repos,
	}
	send(w, r, out, nil)
}

func (s *Server) godocIndexHandler(w http.ResponseWriter, r *http.Request) {
	repos, err := godoc.New().GetIndex()
	send(w, r, repos, err)
}
