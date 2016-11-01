package server

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	exago "github.com/hotolab/exago-svc"
	. "github.com/hotolab/exago-svc/config"
	"github.com/justinas/alice"
)

var (
	ErrInvalidRepository = errors.New("The repository doesn't contain Go code")
	logger               = log.WithField("prefix", "server")
)

type Server struct {
	config exago.Config
}

// New creates a new Server instance.
// Dependencies:
// - DB
// - RepositoryHost
// - Pool
// - Showcaser
func New(options ...exago.Option) *Server {
	var s Server
	for _, option := range options {
		option.Apply(&s.config)
	}
	return &s
}

// ListenAndServe binds the HTTP port and listens for requests.
func (s *Server) ListenAndServe() error {
	repoHandlers := alice.New(
		context.ClearHandler,
		s.recoverHandler,
		s.setLogger,
		s.checkValidRepository,
		// rateLimit,
		// requestLock,
	)
	router := NewRouter()

	router.Get("/project/*repository", repoHandlers.ThenFunc(s.repositoryHandler))
	router.Get("/refresh/*repository", repoHandlers.ThenFunc(s.refreshHandler))
	router.Get("/contents/*repository", repoHandlers.ThenFunc(s.fileHandler))
	router.Get("/cached/*repository", repoHandlers.ThenFunc(s.cachedHandler))
	router.Get("/badge/:type/*repository", repoHandlers.ThenFunc(s.badgeHandler))

	baseHandlers := alice.New(
		context.ClearHandler,
		s.recoverHandler,
	)
	router.Get("/projects/recent", baseHandlers.ThenFunc(s.recentHandler))
	router.Get("/projects/top", baseHandlers.ThenFunc(s.topHandler))
	router.Get("/projects/popular", baseHandlers.ThenFunc(s.popularHandler))
	router.Get("/godocindex", baseHandlers.ThenFunc(s.godocIndexHandler))

	logger.Infof("Listening on port %d", Config.HttpPort)
	return graceful.RunWithErr(fmt.Sprintf("%s:%d", Config.Bind, Config.HttpPort), 10*time.Second, router)
}
