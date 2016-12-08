package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/cors"
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
	processingList processingList
	config         exago.Config
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

	s.processingList = processingList{repos: map[string]bool{}}
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
	r.Use(middleware.Heartbeat("_health"))

	cors := cors.New(cors.Options{
		AllowedOrigins:   Config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Location"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	r.Route("/repos", func(r chi.Router) {
		r.With(s.checkValidData, s.setLogger).Post("/", s.processRepository)
		r.Route("/:repo/branches/:branch/goversions/:goversion", func(r chi.Router) {
			r.Use(s.checkValidData, s.setLogger)
			r.Get("/", s.getRepo)
			r.Get("/files/*path", s.getFile)
			r.Get("/badges/:type", s.getBadge)
		})
	})

	r.Get("/projects/recent", s.getRecent)
	r.Get("/projects/top", s.getTop)
	r.Get("/projects/popular", s.getPopular)

	srv := &graceful.Server{
		Timeout: 20 * time.Second,
		Server:  &http.Server{Addr: fmt.Sprintf("%s:%d", Config.Bind, Config.HttpPort), Handler: chi.ServerBaseContext(r, baseCtx)},
	}
	srv.BeforeShutdown = func() bool {
		logger.Infoln("Shutting down...")
		err := valv.Shutdown(srv.Timeout)
		if err != nil {
			logger.Errorf("Shutdown error %v", err)
		}
		return true
	}

	logger.Infof("Listening on port %d", Config.HttpPort)
	return srv.ListenAndServe()
}
