package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/badge"
	"github.com/exago/svc/config"
	"github.com/exago/svc/datafetcher"
	"github.com/exago/svc/github"
	"github.com/exago/svc/repository"
	"github.com/exago/svc/requestlock"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var (
	errInvalidParameter  = errors.New("Invalid parameter passed")
	errInvalidRepository = errors.New("The repository doesn't contain Go code")
)

func ListenAndServe() {
	repoHandlers := alice.New(
		context.ClearHandler,
		recoverHandler,
		setLogger,
		checkRegistry,
		rateLimit,
		requestLock,
		checkValidRepository,
	)
	router := NewRouter()

	router.Get("/:registry/:owner/:repository/badge/:type/img.svg", repoHandlers.ThenFunc(badgeHandler))
	router.Get("/:registry/:owner/:repository/valid", repoHandlers.ThenFunc(repoValidHandler))
	router.Get("/:registry/:owner/:repository/codestats", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:owner/:repository/imports", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:owner/:repository/lint/:linter", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:owner/:repository/test", repoHandlers.ThenFunc(testHandler))
	router.Get("/:registry/:owner/:repository/rank", repoHandlers.ThenFunc(rankHandler))
	router.Get("/:registry/:owner/:repository/contents/*path", repoHandlers.ThenFunc(fileHandler))

	hp := config.Values.HttpPort
	log.Info("Listening on port " + hp)
	if err := http.ListenAndServe(":"+hp, router); err != nil {
		log.Fatal(err)
	}
}

func badgeHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	lgr := context.Get(r, "logger").(*log.Entry)
	tp := ps.ByName("type")

	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))
	switch tp {
	case "rank":
		title := "exago"
		repo := repository.New(rp, "")
		err, rank := repo.LoadFromDB(), repo.Rank()
		if err != nil {
			lgr.Error(err)
			badge.WriteError(w, title)
			return
		}
		badge.Write(w, title, string(rank), "blue")
	}
}

func rankHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))

	repo := repository.New(rp, "")
	err, rank := repo.LoadFromDB(), repo.Rank()
	send(w, r, rank, err)
}

func lambdaHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)

	lambdaFn := strings.Split(r.URL.String()[1:], "/")[3]
	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))

	var out interface{}
	var err error

	switch lambdaFn {
	case "imports":
		out, err = datafetcher.GetImports(rp)
	case "codestats":
		out, err = datafetcher.GetCodeStats(rp)
	case "lint":
		out, err = datafetcher.GetLintMessages(rp, ps.ByName("linter"))
	}

	send(w, r, out, err)
}

func repoValidHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	owner, repository := ps.ByName("owner"), ps.ByName("repository")
	code, err := checkRepo(r, owner, repository)
	writeData(w, r, code, err)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	content, err := github.GetFileContent(ps.ByName("owner"), ps.ByName("repository"), ps.ByName("path")[1:])
	send(w, r, content, err)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))

	// Add request lock to avoid spawning too many container requests
	requestlock.Add(rp, getIP(r.RemoteAddr), r.URL.Path)

	out, err := datafetcher.GetTestResults(rp)

	// Remove the lock once the data is fetched
	requestlock.Remove(rp, getIP(r.RemoteAddr), r.URL.Path)

	send(w, r, out, err)
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
