package logger

import (
	log "github.com/Sirupsen/logrus"
	. "github.com/exago/svc/config"
)

// SetUp configures the logger.
func SetUp() {
	setLogLevel()
}

func With(repository string, ip string) *log.Entry {
	return log.WithFields(log.Fields{
		"repository": repository,
		"ip":         ip,
	})
}

func setLogLevel() {
	logLevel := Config.LogLevel
	if logLevel == "" {
		log.Info("Log level default to info")
		log.SetLevel(log.InfoLevel)
		return
	}

	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithFields(log.Fields{
			"passed":  logLevel,
			"default": "fatal",
		}).Warn("Log level not valid, fallback to info")
		log.SetLevel(log.InfoLevel)
		return
	}

	log.SetLevel(lvl)
	log.WithFields(log.Fields{
		"level": logLevel,
	}).Debug("Log level successfully set")
}
