package server

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/badge"
	"github.com/exago/svc/github"
	"github.com/exago/svc/repository"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

func repositoryHandler(w http.ResponseWriter, r *http.Request) {
	repo := context.Get(r, "repository").(string)
	rc := repository.NewChecker(repo)

	if rc.Repository.IsCached() {
		err := rc.Repository.Load()
		send(w, r, rc.Repository.FormatOutput(), err)
		return
	}

	// Wait until the data is ready
	rc.Run()
	<-rc.Done
	send(w, r, rc.Output, nil)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	rp := repository.New(context.Get(r, "repository").(string), "")
	rp.ClearCache()
	repositoryHandler(w, r)
}

func badgeHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	lgr := context.Get(r, "logger").(*log.Entry)

	repo := repository.New(ps.ByName("repository")[1:], "")
	isCached := repo.IsCached()
	if !isCached {
		badge.WriteError(w)
		return
	}

	err, rank := repo.Load(), repo.Score.Rank
	if err != nil {
		lgr.Error(err)
		badge.WriteError(w)
		return
	}
	badge.Write(w, string(rank), "blue")
}

func repoValidHandler(w http.ResponseWriter, r *http.Request) {
	owner := context.Get(r, "owner").(string)
	project := context.Get(r, "project").(string)
	code, err := checkRepo(r, owner, project)
	writeData(w, r, code, err)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	owner := context.Get(r, "owner").(string)
	project := context.Get(r, "project").(string)
	path := context.Get(r, "path").(string)
	content, err := github.GetFileContent(owner, project, path)
	send(w, r, content, err)
}

func cachedHandler(w http.ResponseWriter, r *http.Request) {
	repo := context.Get(r, "repository").(string)
	rp := repository.New(repo, "")
	send(w, r, rp.IsCached(), nil)
}

func recentHandler(w http.ResponseWriter, r *http.Request) {
	rp := repository.New(context.Get(r, "repository").(string), "")
	rp.ClearCache()
	repositoryHandler(w, r)
}
