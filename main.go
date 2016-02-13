package main

import (
	"github.com/exago/svc/config"
	"github.com/exago/svc/logger"
	"github.com/exago/svc/redis"
	"github.com/exago/svc/server"
)

func init() {
	config.SetUp()
	logger.SetUp()
	redis.SetUp()
}

func main() {
	server.ListenAndServe()
}
