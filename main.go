package main

import (
	. "github.com/exago/svc/config"
	"github.com/exago/svc/github"
	"github.com/exago/svc/leveldb"
	"github.com/exago/svc/logger"
	"github.com/exago/svc/server"
)

func init() {
	InitializeConfig()
	logger.SetUp()
	github.SetUp()
	leveldb.Init()
}

func main() {
	server.ListenAndServe()
}
