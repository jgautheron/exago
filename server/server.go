package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/badge"
	. "github.com/exago/svc/config"
	"github.com/exago/svc/github"
	"github.com/exago/svc/repository"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var (
	errInvalidParameter  = errors.New("Invalid parameter passed")
	errInvalidRepository = errors.New("The repository doesn't contain Go code")
	errRoutineTimeout    = errors.New("The analysis timed out")
)

func ListenAndServe() {
	repoHandlers := alice.New(
		context.ClearHandler,
		recoverHandler,
		initDB,
		setLogger,
		checkValidRepository,
		// rateLimit,
		// requestLock,
	)
	router := NewRouter()

	router.Get("/project/*repository", repoHandlers.ThenFunc(repositoryHandler))
	router.Get("/badge/*repository", repoHandlers.ThenFunc(badgeHandler))
	router.Get("/valid/*repository", repoHandlers.ThenFunc(repoValidHandler))
	router.Get("/contents/*repository", repoHandlers.ThenFunc(fileHandler))
	router.Get("/cached/*repository", repoHandlers.ThenFunc(cachedHandler))

	log.Infof("Listening on port %d", Config.HttpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", Config.Bind, Config.HttpPort), router))
}

type RepositoryChecker struct {
	repository     *repository.Repository
	types, linters []string
	data           chan interface{}
	dataLoaded     chan bool
	stamped        chan bool
	output         map[string]interface{}
}

func (rc *RepositoryChecker) Stamp() {
	<-rc.dataLoaded

	sc, err := rc.repository.GetScore()
	if err != nil {
		rc.output["score"] = err
	} else {
		rc.output["score"] = sc
	}

	date, err := rc.repository.GetDate()
	if err != nil {
		rc.output["date"] = err
	} else {
		rc.output["date"] = date
	}

	rc.stamped <- true
}

func (rc *RepositoryChecker) RunAll() {
	i := 0
	for _, tp := range rc.types {
		go func(tp string) {
			var (
				out interface{}
				err error
			)

			switch tp {
			case "imports":
				out, err = rc.repository.GetImports()
			case "codestats":
				out, err = rc.repository.GetCodeStats()
			case "testresults":
				out, err = rc.repository.GetTestResults()
			case "lintmessages":
				out, err = rc.repository.GetLintMessages(rc.linters)
			}

			if err != nil {
				rc.data <- err
				return
			}

			rc.data <- out
		}(tp)

		select {
		case out := <-rc.data:
			switch out.(type) {
			case error:
				rc.output[tp] = out.(error)
			default:
				rc.output[tp] = out
				i++
			}
		case <-time.After(time.Minute * 5):
			rc.output[tp] = errRoutineTimeout
		}
	}

	if i == len(rc.types) {
		rc.dataLoaded <- true
	}
}

func repositoryHandler(w http.ResponseWriter, r *http.Request) {
	repo := context.Get(r, "repository").(string)

	rc := &RepositoryChecker{
		repository: repository.New(repo, ""),
		types:      []string{"imports", "codestats", "testresults", "lintmessages"},
		data:       make(chan interface{}, 10),
		dataLoaded: make(chan bool, 1),
		stamped:    make(chan bool, 1),
		linters:    repository.DefaultLinters,
		output:     map[string]interface{}{},
	}
	go rc.RunAll()
	go rc.Stamp()

	<-rc.stamped
	send(w, r, rc.output, nil)
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

// checkRepo ensures that the repository exists on GitHub
// and that it is contains Go code.
func checkRepo(r *http.Request, owner, repo string) (int, error) {
	// Attempt to load the repo
	rp, _, err := github.Repositories.Get(owner, repo)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("Repository %s not found in Github.", repo)
	}

	// Check if the repo contains Go code
	if !strings.Contains(*rp.Language, "Go") {
		return http.StatusNotAcceptable, errInvalidRepository
	}

	return http.StatusOK, nil
}
