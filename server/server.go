package server

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	log "github.com/Sirupsen/logrus"
	. "github.com/exago/svc/config"
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

	log.Infof("Listening on port %d", Config.HttpPort)
	graceful.Run(fmt.Sprintf("%s:%d", Config.Bind, Config.HttpPort), 10*time.Second, router)
}
