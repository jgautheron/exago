package server

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	. "github.com/hotolab/exago-svc/config"
	"github.com/justinas/alice"
)

var (
	ErrInvalidRepository = errors.New("The repository doesn't contain Go code")
	logger               = log.WithField("prefix", "server")
)

// ListenAndServe binds the HTTP port and listens for requests.
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
	router.Get("/contents/*repository", repoHandlers.ThenFunc(fileHandler))
	router.Get("/cached/*repository", repoHandlers.ThenFunc(cachedHandler))
	router.Get("/badge/:type/*repository", repoHandlers.ThenFunc(badgeHandler))

	baseHandlers := alice.New(
		context.ClearHandler,
		recoverHandler,
	)
	router.Get("/projects/recent", baseHandlers.ThenFunc(recentHandler))
	router.Get("/projects/top", baseHandlers.ThenFunc(topHandler))
	router.Get("/projects/popular", baseHandlers.ThenFunc(popularHandler))
	router.Get("/godocindex", baseHandlers.ThenFunc(godocIndexHandler))

	logger.Infof("Listening on port %d", Config.HttpPort)
	logger.Fatal(graceful.RunWithErr(fmt.Sprintf("%s:%d", Config.Bind, Config.HttpPort), 10*time.Second, router))
}
