package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
	LogFormat string `envconfig:"LOG_FORMAT" default:"json"`
}

type HTTPConfig struct {
	HTTPPort       int      `envconfig:"HTTP_PORT" default:"8080"`
	HTTPBind       string   `envconfig:"HTTP_BIND" default:"0.0.0.0"`
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" default:"*"`
}

type GitHubConfig struct {
	GithubAccessTokens []string `envconfig:"GITHUB_ACCESS_TOKENS" required:"true"`
}

type GoogleCloudConfig struct {
	GoogleProjectID             string `envconfig:"GCLOUD_PROJECT_ID" required:"true"`
	GooglePubSubTopicRepository string `envconfig:"GCLOUD_PUBSUB_TOPIC_REPOSITORY" required:"true"`
}

func InitializeConfig(target interface{}) {
	if err := envconfig.Process("", target); err != nil {
		logrus.Fatal(err)
	}
}

// InitializeLogging sets logrus log level and formatting style.
func InitializeLogging(logLevel, logFormat string) {
	switch strings.ToLower(logFormat) {
	case "text":
		logrus.SetFormatter(new(logrus.TextFormatter))
	default:
		logrus.SetFormatter(new(logrus.JSONFormatter))
	}

	// If log level cannot be resolved, exit gracefully
	if logLevel == "" {
		logrus.Warning("Log level could not be resolved, fallback to fatal")
		logrus.SetLevel(logrus.FatalLevel)
		return
	}

	// Parse level from string
	lvl, err := logrus.ParseLevel(logLevel)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"passed":  logLevel,
			"default": "fatal",
		}).Warn("Log level is not valid, fallback to default level")
		logrus.SetLevel(logrus.FatalLevel)
		return
	}

	logrus.SetLevel(lvl)
	logrus.WithFields(logrus.Fields{
		"lvl": logLevel,
	}).Debug("Log level successfully set")
}
