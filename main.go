package main

import (
	"github.com/jgautheron/exago-service/config"
	"github.com/jgautheron/exago-service/logger"
	"github.com/jgautheron/exago-service/redis"
	"github.com/jgautheron/exago-service/server"
)

func init() {
	config.SetUp()
	logger.SetUp()
	redis.SetUp()
}

func main() {
	server.ListenAndServe()
}
