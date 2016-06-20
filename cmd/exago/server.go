package main

import (
	"github.com/codegangsta/cli"
	"github.com/exago/svc/github"
	"github.com/exago/svc/server"
)

// ServerCommand starts the HTTP server.
func ServerCommand() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start the HTTP server",
		Action: func(ctx *cli.Context) {
			github.SetUp()
			server.ListenAndServe()
		},
	}
}
