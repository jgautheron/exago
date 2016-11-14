package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/hotolab/exago-svc/badge"
	"github.com/hotolab/exago-svc/repository/model"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) repositoryHandler(w http.ResponseWriter, r *http.Request) {
	handleAction := func(rp model.Record, err error) {
		if err != nil {
			send(w, r, nil, err)
			return
		}
		go s.config.Showcaser.Process(rp)
		send(w, r, rp, nil)
	}

	repo := context.Get(r, "repository").(string)
	branch := ""

	if s.config.RepositoryLoader.IsCached(repo, branch) {
		rp, err := s.config.RepositoryLoader.Load(repo, branch)
		handleAction(rp, err)
		return
	}

	rp, err := s.config.Pool.PushSync(repo)
	handleAction(rp, err)
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

	var title string
	var value string
	var score = rp.GetScore().Value

	pr := rp.GetProjectRunner()
	th := pr.Thirdparties.Data
	cs := pr.CodeStats.Data

	switch tp := ps.ByName("type"); tp {
	case "cov":
		title = "coverage"
		value = fmt.Sprintf("%.2f%%", pr.GetMeanCodeCov())
	case "duration":
		title = "tests duration"
		value = fmt.Sprintf("%.2fs", pr.GetMeanTestDuration())
	case "tests":
		title = "tests"
		value = fmt.Sprintf("%d", cs["test"])
	case "thirdparties":
		title = "3rd parties"
		value = fmt.Sprintf("%d", len(th))
	case "loc":
		title = "LOC"
		value = fmt.Sprintf("%d", cs["loc"])
	case "rank":
		fallthrough
	default:
		value = rp.GetRank()
	}

	badge.Write(w, title, value, score)
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
		send(w, r, rp, err)
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

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	writeData(w, r, http.StatusOK, `{"alive": true}`)
}
