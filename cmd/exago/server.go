package main

import (
	"github.com/exago/svc/github"
	"github.com/exago/svc/server"
	"github.com/exago/svc/showcaser"
	"github.com/urfave/cli"
)

// ServerCommand starts the HTTP server.
func ServerCommand() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start the HTTP server",
		Action: func(ctx *cli.Context) {
			github.Init()
			showcaser.Init()
			server.ListenAndServe()
		},
	}
}
