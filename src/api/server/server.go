package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jgautheron/exago/src/api/database/firestore"
	"github.com/jgautheron/exago/src/api/eventpub"
)

var ErrInvalidRepository = errors.New("The repository doesn't contain Go code")

type Server struct {
	db  *firestore.Firestore
	evp *eventpub.EventPub
}

func New() (*Server, error) {
	ctx := context.Background()

	db, err := firestore.NewFromConfig(ctx, &Config.GoogleCloudConfig)
	if err != nil {
		return nil, err
	}

	evp, err := eventpub.NewFromConfig(ctx, &Config.GoogleCloudConfig)
	if err != nil {
		return nil, err
	}

	return &Server{db: db, evp: evp}, nil
}

// ListenAndServe binds the HTTP port and listens for requests.
func (s *Server) ListenAndServe() error {
	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins:   Config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "PATCH", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cache-Control"},
		ExposedHeaders:   []string{"Location"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	r.Get("/project/{repo}", s.testHandler)
	r.Get("/refresh/{repo}", s.testHandler)
	r.Get("/contents/{repo}", s.testHandler)
	r.Get("/badge/{type}/{repo}", s.testHandler)

	r.Get("/projects/recent", s.testHandler)
	r.Get("/projects/top", s.testHandler)
	r.Get("/projects/popular", s.testHandler)

	return http.ListenAndServe(
		fmt.Sprintf("%s:%d", Config.HTTPBind, Config.HTTPPort),
		r,
	)
}
