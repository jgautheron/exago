package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/tylerb/graceful"

	"net/http"

	log "github.com/Sirupsen/logrus"
	exago "github.com/hotolab/exago-svc"
	. "github.com/hotolab/exago-svc/config"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/valve"
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
	// Our graceful valve shut-off package to manage code preemption and
	// shutdown signaling.
	valv := valve.New()
	baseCtx := valv.Context()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CloseNotify)

	r.Route("/repos", func(r chi.Router) {
		r.Use(s.checkValidRepository)
		r.Post("/", s.repositoryHandler)
		r.Route("/:host|:owner|:repo/branches/:branch/goversion/:goversion", func(r chi.Router) {
			r.Get("/", s.cachedHandler)
			r.Get("/file/*path", s.fileHandler)
			r.Get("/badges/:type", s.badgeHandler)
		})
	})

	r.Get("/_health", s.healthHandler)

	// router.Get("/project/*repository", repoHandlers.ThenFunc(s.repositoryHandler))
	// router.Get("/refresh/*repository", repoHandlers.ThenFunc(s.refreshHandler))
	// router.Get("/contents/*repository", repoHandlers.ThenFunc(s.fileHandler))
	// router.Get("/cached/*repository", repoHandlers.ThenFunc(s.cachedHandler))
	// router.Get("/badge/:type/*repository", repoHandlers.ThenFunc(s.badgeHandler))

	// baseHandlers := alice.New(
	// 	context.ClearHandler,
	// 	s.recoverHandler,
	// )
	// router.Get("/projects/recent", baseHandlers.ThenFunc(s.recentHandler))
	// router.Get("/projects/top", baseHandlers.ThenFunc(s.topHandler))
	// router.Get("/projects/popular", baseHandlers.ThenFunc(s.popularHandler))

	// r.Get("/_health", baseHandlers.ThenFunc(s.healthHandler))

	srv := &graceful.Server{
		Timeout: 20 * time.Second,
		Server:  &http.Server{Addr: fmt.Sprintf("%s:%d", Config.Bind, Config.HttpPort), Handler: chi.ServerBaseContext(r, baseCtx)},
	}
	srv.BeforeShutdown = func() bool {
		fmt.Println("shutting down..")
		err := valv.Shutdown(srv.Timeout)
		if err != nil {
			fmt.Println("Shutdown error -", err)
		}

		// the app code has stopped here now, and so this would be a good place
		// to close up any db and other service connections, etc.
		return true
	}
	logger.Infof("Listening on port %d", Config.HttpPort)
	return srv.ListenAndServe()
}
