package logger

import (
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/config"
)

// SetUp configures logrus to send logs to Papertrail if
// we are in a production environment.
func SetUp() {
	setLogLevel()
}

func With(repository string, ip string) *log.Entry {
	// Retrieve the user IP
	if ip != "" {
		ip, _, err := net.SplitHostPort(ip)
		if err != nil {
			log.Error("userip: %q is not IP:port", ip)
			return &log.Entry{}
		}
	}

	return log.WithFields(log.Fields{
		"repository": repository,
		"ip":         ip,
	})
}

func setLogLevel() {
	logLevel := config.Get("LOG_LEVEL")

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
