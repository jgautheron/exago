package config

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/envconfig"
)

var Config cfg

type cfg struct {
	GithubAccessToken  string `envconfig:"GITHUB_ACCESS_TOKEN" required:"true"`
	AwsRegion          string `envconfig:"AWS_REGION" default:"eu-west-1"`
	AwsAccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID" required:"true"`
	AwsSecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY" required:"true"`
	HttpPort           int    `envconfig:"HTTP_PORT" default:"8080"`
	Bind               string `envconfig:"BIND" default:"127.0.0.1"`
	DatabasePath       string `envconfig:"DATABASE_PATH" default:"./exago.db"`
	AllowOrigin        string `envconfig:"ALLOW_ORIGIN" default:"*"`
	LogLevel           string `envconfig:"LOG_LEVEL" default:"info"`

	ShowcaserPopularRebuildInterval time.Duration `envconfig:"SHOWCASER_POPULAR_REBUILD_INTERVAL" default:"1m"`

	PoolSize int `envconfig:"POOL_SIZE" default:"20"`
}

func InitializeConfig() {
	if err := envconfig.Process("", &Config); err != nil {
		log.Fatal(err)
	}
}
