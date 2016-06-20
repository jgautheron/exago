package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	. "github.com/exago/svc/config"
	"github.com/exago/svc/logger"
)

var (
	App *cli.App
)

// Initialize commandline app.
func init() {
	App = cli.NewApp()

	// For fancy output on console
	App.Name = "exago godoc indexer"
	App.Usage = `Indexes the godoc index in DB`
	App.Author = "Jonathan Gautheron"

	// Version is injected at build-time
	App.Version = ""

	InitializeConfig()
	logger.SetUp()
}

func main() {
	AddCommands()
	if err := App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// AddCommands adds child commands to the root command Cmd.
func AddCommands() {
	AddCommand(IndexCommand())
}

// AddCommand adds a child command.
func AddCommand(cmd cli.Command) {
	App.Commands = append(App.Commands, cmd)
}
