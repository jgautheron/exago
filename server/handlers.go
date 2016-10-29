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
		err, rank, score := repo.Load(), repo.GetRank(), repo.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, "")
			return
		}

		badge.Write(w, "", string(rank), score)
	case "cov":
		title := "coverage"
		err, cov, score := repo.Load(), repo.GetProjectRunner().GetMeanCodeCov(), repo.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%.2f%%", cov), score)
	case "duration":
		title := "tests duration"
		err, avg, score := repo.Load(), repo.GetProjectRunner().GetMeanTestDuration(), repo.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%.2fs", avg), score)
	case "tests":
		title := "tests"
		err, tests, score := repo.Load(), repo.GetCodeStats()["Test"], repo.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", tests), score)
	case "thirdparties":
		title := "3rd parties"
		err, thirdParties, score := repo.Load(), len(repo.GetProjectRunner().Thirdparties.Data), repo.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", thirdParties), score)
	case "loc":
		title := "LOC"
		err, thirdParties, score := repo.Load(), repo.GetCodeStats()["LOC"], repo.GetScore().Value
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", thirdParties), score)
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
