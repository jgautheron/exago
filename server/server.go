package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/config"
	"github.com/exago/svc/datafetcher"
	"github.com/exago/svc/rank"
	"github.com/exago/svc/requestlock"
	"github.com/google/go-github/github"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var (
	errInvalidParameter  = errors.New("Invalid parameter passed")
	errInvalidRepository = errors.New("The repository doesn't contain Go code")
	errBadgeError        = errors.New("The badge cannot be fetched")
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
		// checkCache,
	)
	router := NewRouter()

	// router.Get("/:registry/:owner/:repository/badge/:type/img.svg", repoHandlers.ThenFunc(badgeHandler))
	router.Get("/:registry/:owner/:repository/valid", repoHandlers.ThenFunc(repoValidHandler))
	router.Get("/:registry/:owner/:repository/loc", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:owner/:repository/imports", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:owner/:repository/lint/:linter", repoHandlers.ThenFunc(lambdaHandler))
	router.Get("/:registry/:owner/:repository/test", repoHandlers.ThenFunc(testHandler))
	router.Get("/:registry/:owner/:repository/rank", repoHandlers.ThenFunc(rankHandler))
	router.Get("/:registry/:owner/:repository/contents/*path", repoHandlers.ThenFunc(fileHandler))

	hp := config.Get("HttpPort")
	log.Info("Listening on port " + hp)
	if err := http.ListenAndServe(":"+hp, router); err != nil {
		log.Fatal(err)
	}
}

// func badgeHandler(w http.ResponseWriter, r *http.Request) {
// 	ps := context.Get(r, "params").(httprouter.Params)
// 	tp := ps.ByName("type")

// 	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))
// 	switch tp {
// 	case "imports":
// 		out, err := datafetcher.GetImports(rp)
// 		if err != nil {
// 			log.Error(err)
// 			badge.WriteError(w, tp)
// 			return
// 		}
// 		var data []string
// 		err = json.Unmarshal(*out, &data)
// 		if err != nil {
// 			log.Error(err)
// 			badge.WriteError(w, tp)
// 			return
// 		}
// 		ln := strconv.Itoa(len(data))
// 		badge.Write(w, "imports", ln, "blue")
// 	case "rank":
// 		title := "exago"

// 		rk := rank.New()
// 		rk.SetRepository(rp)
// 		val, err := rk.GetRankFromCache()
// 		if err != nil {
// 			log.Error(err)
// 			badge.WriteError(w, title)
// 			return
// 		}
// 		badge.Write(w, title, val, "blue")
// 	}
// }

func rankHandler(w http.ResponseWriter, r *http.Request) {
	ps := context.Get(r, "params").(httprouter.Params)
	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))

	rk := rank.New()
	rk.SetRepository(rp)
	rank, err := rk.GetScore()
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
	case "loc":
		out, err = datafetcher.GetCodeStats(rp)
	case "lint":
		out, err = datafetcher.GetLint(rp, ps.ByName("linter"))
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
	lgr := context.Get(r, "logger").(*log.Entry)

	rp := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))
	out, err := datafetcher.GetFileContents(lgr, rp, ps.ByName("path")[1:])
	send(w, r, string(out), err)
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
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Get("GithubAccessToken")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	// Attempt to load the repo
	rp, _, err := client.Repositories.Get(owner, repo)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("Repository %s not found in Github.", repo)
	}

	// Check if the repo contains Go code
	if !strings.Contains(*rp.Language, "Go") {
		return http.StatusNotAcceptable, errInvalidRepository
	}

	return http.StatusOK, nil
}
