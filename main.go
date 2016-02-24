package main

import (
	"github.com/exago/svc/config"
	"github.com/exago/svc/logger"
	"github.com/exago/svc/server"
)

func init() {
	config.SetUp()
	logger.SetUp()
}

func main() {
	server.ListenAndServe()
}
