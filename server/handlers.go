package server

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/badge"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

const (
	WebsiteURL = "https://exago.io"
)

func (s *Server) processRepository(w http.ResponseWriter, r *http.Request) {
	repo := r.Context().Value("repo").(string)
	branch := r.Context().Value("branch").(string)
	goversion := r.Context().Value("goversion").(string)

	// Avoid double processing
	if s.processingList.exists(repo, branch, goversion) {
		render.Status(r, http.StatusTooManyRequests)
		render.JSON(w, r, http.StatusText(http.StatusTooManyRequests))
		return
	}

	s.processingList.add(repo, branch, goversion)
	defer s.processingList.remove(repo, branch, goversion)

	_, err := s.config.Pool.PushSync(repo, branch, goversion)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf(
		"%s/repos/%s/branches/%s/goversions/%s",
		WebsiteURL,
		strings.Replace(repo, "/", "|", 2),
		branch,
		goversion,
	))
	render.NoContent(w, r)
}

func (s *Server) getBadge(w http.ResponseWriter, r *http.Request) {
	repo := r.Context().Value("repo").(string)
	branch := r.Context().Value("branch").(string)
	goversion := r.Context().Value("goversion").(string)
	lgr := r.Context().Value("lgr").(*log.Entry)

	if !s.config.RepositoryLoader.IsCached(repo, branch, goversion) {
		badge.WriteError(w, "")
		return
	}

	rp, err := s.config.RepositoryLoader.Load(repo, branch, goversion)
	if err != nil {
		lgr.Error(err)
		badge.WriteError(w, "")
		return
	}

	var title, value string
	var score = rp.GetScore().Value

	pr := rp.GetProjectRunner()
	th := pr.Thirdparties.Data
	cs := pr.CodeStats.Data

	switch tp := chi.URLParam(r, "type"); tp {
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

func (s *Server) getRepo(w http.ResponseWriter, r *http.Request) {
	repo := r.Context().Value("repo").(string)
	branch := r.Context().Value("branch").(string)
	goversion := r.Context().Value("goversion").(string)

	rp, err := s.config.RepositoryLoader.Load(repo, branch, goversion)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, http.StatusText(http.StatusNotFound))
		return
	}
	go s.config.Showcaser.Process(rp)
	render.JSON(w, r, rp)
}

func (s *Server) getFile(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value("owner").(string)
	name := r.Context().Value("name").(string)
	// branch := r.Context().Value("branch").(string)
	path := chi.URLParam(r, "path")
	content, err := s.config.RepositoryHost.GetFileContent(owner, name, path)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}
	render.PlainText(w, r, content)
}

func (s *Server) getProjects(w http.ResponseWriter, r *http.Request) {
	tp := r.URL.Query().Get("type")
	var repos []interface{}
	switch tp {
	case "recent":
		repos = s.config.Showcaser.GetRecentRepositories()
	case "top":
		repos = s.config.Showcaser.GetTopRankedRepositories()
	case "popular":
		repos = s.config.Showcaser.GetPopularRepositories()
	}
	list := map[string]interface{}{
		"type":         tp,
		"repositories": repos,
	}
	render.JSON(w, r, list)
}

func (s *Server) getBranches(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value("owner").(string)
	name := r.Context().Value("name").(string)
	res, err := s.config.RepositoryHost.Get(owner, name)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}
	render.JSON(w, r, res["branches"])
}