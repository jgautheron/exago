package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/badge"
	"github.com/exago/svc/github"
	"github.com/exago/svc/godoc"
	"github.com/exago/svc/repository"
	"github.com/exago/svc/repository/processor"
	"github.com/exago/svc/showcaser"
	"github.com/exago/svc/taskrunner/lambda"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

func repositoryHandler(w http.ResponseWriter, r *http.Request) {
	repo := context.Get(r, "repository").(string)
	rc := processor.NewChecker(repo, lambda.Runner{Repository: repo})

	if rc.Repository.IsCached() {
		err := rc.Repository.Load()
		send(w, r, rc.Repository.AsMap(), err)
		return
	}

	rc.Run()

	// Wait until the data is ready
	select {
	case <-rc.Done:
		send(w, r, rc.Output, nil)
	case <-rc.Aborted:
		send(w, r, rc.Output, nil)
	}
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
		badge.WriteError(w, "")
		return
	}

	switch tp := ps.ByName("type"); tp {
	case "rank":
		err, rank := repo.Load(), repo.Score.Rank
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, "")
			return
		}
		badge.Write(w, "", string(rank), "blue")
	case "cov":
		title := "coverage"
		err, cov := repo.Load(), repo.TestResults.GetAvgCodeCov()
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%.2f%%", cov), "blue")
	case "duration":
		title := "tests duration"
		err, avg := repo.Load(), repo.TestResults.GetAvgTestDuration()
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%.2fs", avg), "blue")
	case "tests":
		title := "tests"
		err, tests := repo.Load(), repo.CodeStats["Test"]
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", tests), "blue")
	case "thirdparties":
		title := "3rd parties"
		err, thirdParties := repo.Load(), len(repo.Imports)
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", thirdParties), "blue")
	case "loc":
		title := "LOC"
		err, thirdParties := repo.Load(), repo.CodeStats["LOC"]
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, fmt.Sprintf("%d", thirdParties), "blue")
	}
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
	repos := showcaser.GetRecentRepositories()
	out := map[string]interface{}{
		"type":         "recent",
		"repositories": repos,
	}
	send(w, r, out, nil)
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	repos := showcaser.GetTopRankedRepositories()
	out := map[string]interface{}{
		"type":         "top",
		"repositories": repos,
	}
	send(w, r, out, nil)
}

func popularHandler(w http.ResponseWriter, r *http.Request) {
	repos := showcaser.GetPopularRepositories()
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
