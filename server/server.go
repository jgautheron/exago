package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	log "github.com/Sirupsen/logrus"
	. "github.com/exago/svc/config"
	"github.com/exago/svc/github"
	"github.com/gorilla/context"
	"github.com/justinas/alice"
)

var (
	ErrInvalidRepository = errors.New("The repository doesn't contain Go code")
)

func ListenAndServe() {
	repoHandlers := alice.New(
		context.ClearHandler,
		recoverHandler,
		setLogger,
		checkValidRepository,
		// rateLimit,
		// requestLock,
	)
	router := NewRouter()

	router.Get("/project/*repository", repoHandlers.ThenFunc(repositoryHandler))
	router.Get("/refresh/*repository", repoHandlers.ThenFunc(refreshHandler))
	router.Get("/valid/*repository", repoHandlers.ThenFunc(repoValidHandler))
	router.Get("/contents/*repository", repoHandlers.ThenFunc(fileHandler))
	router.Get("/cached/*repository", repoHandlers.ThenFunc(cachedHandler))
	router.Get("/badge/:type/*repository", repoHandlers.ThenFunc(badgeHandler))

	projectHandlers := alice.New(
		context.ClearHandler,
		recoverHandler,
	)
	router.Get("/projects/recent", projectHandlers.ThenFunc(recentHandler))
	router.Get("/projects/top", projectHandlers.ThenFunc(topHandler))
	router.Get("/projects/popular", projectHandlers.ThenFunc(popularHandler))

	log.Infof("Listening on port %d", Config.HttpPort)
	graceful.Run(fmt.Sprintf("%s:%d", Config.Bind, Config.HttpPort), 10*time.Second, router)
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
		return http.StatusNotAcceptable, ErrInvalidRepository
	}

	return http.StatusOK, nil
}
