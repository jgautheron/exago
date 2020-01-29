package server

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jgautheron/exago/internal/eventpub"
	"github.com/sirupsen/logrus"
)

func (s Server) testHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}

func (s Server) processRepository(w http.ResponseWriter, r *http.Request) {
	branch := chi.URLParam(r, "branch")
	repository := chi.URLParam(r, "*")
	goVersion := chi.URLParam(r, "goVersion")

	// Simple input validation
	rg := regexp.MustCompile(`^([a-z0-9\.\-_/]+)$`)
	if !rg.MatchString(branch) || !rg.MatchString(repository) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if match, _ := regexp.MatchString(`^([0-9\.]+)$`, goVersion); !match {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ev := eventpub.RepositoryAddedEvent{
		Branch:     branch,
		Repository: repository,
		GoVersion:  goVersion,
	}
	if err := s.evp.RepositoryAdded(&ev); err != nil {
		logrus.Errorf("Could not enqueue repo: %#v", ev)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
}

//func (s *Server) badgeHandler(w http.ResponseWriter, r *http.Request) {
//	ps := context.Get(r, "params").(httprouter.Params)
//	lgr := context.Get(r, "lgr").(*log.Entry)
//
//	repo := context.Get(r, "repository").(string)
//	branch := ""
//
//	if !s.config.RepositoryLoader.IsCached(repo, branch) {
//		badge.WriteError(w, "")
//		return
//	}
//
//	rp, err := s.config.RepositoryLoader.Load(repo, branch)
//	if err != nil {
//		lgr.Error(err)
//		badge.WriteError(w, "")
//		return
//	}
//
//	var title string
//	var value string
//	var score = rp.GetScore().Value
//
//	pr := rp.GetProjectRunner()
//	th := pr.Thirdparties.Data
//	cs := pr.CodeStats.Data
//
//	switch tp := ps.ByName("type"); tp {
//	case "cov":
//		title = "coverage"
//		value = fmt.Sprintf("%.2f%%", pr.GetMeanCodeCov())
//	case "duration":
//		title = "tests duration"
//		value = fmt.Sprintf("%.2fs", pr.GetMeanTestDuration())
//	case "tests":
//		title = "tests"
//		value = fmt.Sprintf("%d", cs["test"])
//	case "thirdparties":
//		title = "3rd parties"
//		value = fmt.Sprintf("%d", len(th))
//	case "loc":
//		title = "LOC"
//		value = fmt.Sprintf("%d", cs["loc"])
//	case "rank":
//		fallthrough
//	default:
//		value = rp.GetRank()
//	}
//
//	badge.Write(w, title, value, score)
//}
//
//func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request) {
//	owner := context.Get(r, "owner").(string)
//	project := context.Get(r, "project").(string)
//	path := context.Get(r, "path").(string)
//	content, err := s.config.RepositoryHost.GetFileContent(owner, project, path)
//	send(w, r, content, err)
//}
