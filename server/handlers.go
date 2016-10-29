package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/hotolab/exago-svc/badge"
	"github.com/hotolab/exago-svc/github"
	"github.com/hotolab/exago-svc/godoc"
	"github.com/hotolab/exago-svc/pool"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/showcaser"
	"github.com/julienschmidt/httprouter"
)

func repositoryHandler(w http.ResponseWriter, r *http.Request) {
	repo := context.Get(r, "repository").(string)
	rp := repository.New(repo, "")
	if rp.IsCached() {
		err := rp.Load()
		if !rp.HasError() {
			go showcaser.GetInstance().Process(rp)
		}
		send(w, r, rp.GetData(), err)
		return
	}

	rp, err := pool.GetInstance().PushSync(repo)
	if err == nil && !rp.HasError() {
		go showcaser.GetInstance().Process(rp)
	}
	send(w, r, rp.GetData(), err)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	rp := repository.New(context.Get(r, "repository").(string), "")
	rp.ClearCache()
	repositoryHandler(w, r)
}

func badgeHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	lgr := context.Get(r, "lgr").(*log.Entry)

	repo := repository.New(ps.ByName("repository")[1:], "")
	if !repo.IsCached() {
		badge.WriteError(w, "")
		return
	}

	switch tp := ps.ByName("type"); tp {
	case "rank":
		err, rank := repo.Load(), repo.GetRank()
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, "")
			return
		}
		badge.Write(w, "", string(rank), "blue")
	case "cov":
		title := "coverage"
		err, cov := repo.Load(), repo.GetProjectRunner().GetMeanCodeCov()
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%.2f%%", cov), "blue")
	case "duration":
		title := "tests duration"
		err, avg := repo.Load(), repo.GetProjectRunner().GetMeanTestDuration()
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%.2fs", avg), "blue")
	case "tests":
		title := "tests"
		err, tests := repo.Load(), repo.GetCodeStats()["Test"]
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", tests), "blue")
	case "thirdparties":
		title := "3rd parties"
		err, thirdParties := repo.Load(), len(repo.GetProjectRunner().Thirdparties.Data)
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", thirdParties), "blue")
	case "loc":
		title := "LOC"
		err, thirdParties := repo.Load(), repo.GetCodeStats()["LOC"]
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", thirdParties), "blue")
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	owner := context.Get(r, "owner").(string)
	project := context.Get(r, "project").(string)
	path := context.Get(r, "path").(string)
	content, err := github.GetInstance().GetFileContent(owner, project, path)
	send(w, r, content, err)
}

func cachedHandler(w http.ResponseWriter, r *http.Request) {
	repo := context.Get(r, "repository").(string)
	rp := repository.New(repo, "")
	if rp.IsCached() {
		err := rp.Load()
		send(w, r, rp.GetData(), err)
		return
	}
	send(w, r, false, nil)
}

func recentHandler(w http.ResponseWriter, r *http.Request) {
	repos := showcaser.GetInstance().GetRecentRepositories()
	out := map[string]interface{}{
		"type":         "recent",
		"repositories": repos,
	}
	send(w, r, out, nil)
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	repos := showcaser.GetInstance().GetTopRankedRepositories()
	out := map[string]interface{}{
		"type":         "top",
		"repositories": repos,
	}
	send(w, r, out, nil)
}

func popularHandler(w http.ResponseWriter, r *http.Request) {
	repos := showcaser.GetInstance().GetPopularRepositories()
	out := map[string]interface{}{
		"type":         "popular",
		"repositories": repos,
	}
	send(w, r, out, nil)
}

func godocIndexHandler(w http.ResponseWriter, r *http.Request) {
	repos, err := godoc.New().GetIndex()
	send(w, r, repos, err)
}
