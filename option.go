// Package exago exposes first-class functions that are used in the entire application to pass
// around domain instances.
// Greatly inspired from a talk by Dave Cheney: https://youtu.be/5buaPyJ0XeQ
package exago

import "github.com/hotolab/exago-svc/repository/model"

// Option defines the behavior of the configuration settings passed to constructors.
type Option interface {
	Apply(*Config)
}

// Config contains the settings that can be passed to a constructor compatible
// with first-class functions.
type Config struct {
	DB                  model.Database
	Pool                model.Pool
	RepositoryProcessor func(value interface{}) interface{}
	RepositoryHost      model.RepositoryHost
	RepositoryLoader    model.RepositoryLoader
	Showcaser           model.Promoter
}

type database struct {
	db model.Database
}

func (d *database) Apply(cfg *Config) {
	cfg.DB = d.db
}

func WithDatabase(db model.Database) Option {
	return &database{db}
}

type repositoryHost struct {
	rh model.RepositoryHost
}

func (r *repositoryHost) Apply(cfg *Config) {
	cfg.RepositoryHost = r.rh
}

func WithRepositoryHost(rh model.RepositoryHost) Option {
	return &repositoryHost{rh}
}

type repositoryLoader struct {
	rl model.RepositoryLoader
}

func (r *repositoryLoader) Apply(cfg *Config) {
	cfg.RepositoryLoader = r.rl
}

func WithRepositoryLoader(rl model.RepositoryLoader) Option {
	return &repositoryLoader{rl}
}

type processor struct {
	pr func(value interface{}) interface{}
}

func (r *processor) Apply(cfg *Config) {
	cfg.RepositoryProcessor = r.pr
}

func WithProcessor(pr func(value interface{}) interface{}) Option {
	return &processor{pr}
}

type pool struct {
	pl model.Pool
}

func (p *pool) Apply(cfg *Config) {
	cfg.Pool = p.pl
}

func WithPool(pl model.Pool) Option {
	return &pool{pl}
}

type showcaser struct {
	sh model.Promoter
}

func (s *showcaser) Apply(cfg *Config) {
	cfg.Showcaser = s.sh
}

func WithShowcaser(sh model.Promoter) Option {
	return &showcaser{sh}
}
