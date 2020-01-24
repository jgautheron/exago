package server

import (
	"github.com/jgautheron/exago/src/api/internal/config"
)

var Config Cfg

type Cfg struct {
	config.LogConfig
	config.HTTPConfig
	config.GoogleCloudConfig
}

func InitializeConfig() {
	config.InitializeConfig(&Config)
	config.InitializeLogging(Config.LogLevel, Config.LogFormat)
}
