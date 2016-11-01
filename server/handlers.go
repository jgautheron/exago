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
	rp := s.config.RepositoryLoader.Load(repo, "")
	if rp.IsCached() {
		err := rp.Load()
		if !rp.HasError() {
			go s.config.Showcaser.Process(rp)
		}
		send(w, r, rp.GetData(), err)
		return
	}

	rp, err := s.config.Pool.PushSync(repo)
	if err == nil && !rp.HasError() {
		go s.config.Showcaser.Process(rp)
	}
	send(w, r, rp.GetData(), err)
}

func (s *Server) refreshHandler(w http.ResponseWriter, r *http.Request) {
	rp := s.config.RepositoryLoader.Load(
		context.Get(r, "repository").(string),
		"",
	)
	rp.ClearCache()
	s.repositoryHandler(w, r)
}

func (s *Server) badgeHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	lgr := context.Get(r, "lgr").(*log.Entry)

	rp := s.config.RepositoryLoader.Load(
		ps.ByName("repository")[1:],
		"",
	)
	if !rp.IsCached() {
		badge.WriteError(w, "")
		return
	}

	switch tp := ps.ByName("type"); tp {
	case "rank":
		err, rank, score := rp.Load(), rp.GetRank(), rp.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, "")
			return
		}

		badge.Write(w, "", string(rank), score)
	case "cov":
		title := "coverage"
		err, cov, score := rp.Load(), rp.GetProjectRunner().GetMeanCodeCov(), rp.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%.2f%%", cov), score)
	case "duration":
		title := "tests duration"
		err, avg, score := rp.Load(), rp.GetProjectRunner().GetMeanTestDuration(), rp.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%.2fs", avg), score)
	case "tests":
		title := "tests"
		err, tests, score := rp.Load(), rp.GetCodeStats()["Test"], rp.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", tests), score)
	case "thirdparties":
		title := "3rd parties"
		err, thirdParties, score := rp.Load(), len(rp.GetProjectRunner().Thirdparties.Data), rp.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", thirdParties), score)
	case "loc":
		title := "LOC"
		err, thirdParties, score := rp.Load(), rp.GetCodeStats()["LOC"], rp.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
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
	rp := s.config.RepositoryLoader.Load(repo, "")
	if rp.IsCached() {
		err := rp.Load()
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
