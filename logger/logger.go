package logger

import (
	"net"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/jgautheron/exago-service/config"
	"gopkg.in/polds/logrus-papertrail-hook.v2"
)

// SetUp configures logrus to send logs to Papertrail if
// we are in a production environment.
func SetUp() error {
	setLogLevel()

	// Send logs to Papertrail if we are in a production environment.
	if config.Get("Env") == "prod" {
		ptp, _ := strconv.Atoi(config.Get("PapertrailPort"))
		hook, err := logrus_papertrail.NewPapertrailHook(&logrus_papertrail.Hook{
			Host:     config.Get("PapertrailURL"),
			Port:     ptp,
			Hostname: config.Get("PapertrailHost"),
			Appname:  config.Get("PapertrailApp"),
		})
		if err != nil {
			return err
		}
		log.AddHook(hook)

		log.WithFields(log.Fields{
			"driver": "papertrail",
		}).Debug("Logging driver successfully set")
	}

	return nil
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
