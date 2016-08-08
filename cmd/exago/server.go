package main

import (
	"github.com/hotolab/exago-svc/server"
	"github.com/urfave/cli"
)

// ServerCommand starts the HTTP server.
func ServerCommand() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start the HTTP server",
		Action: func(c *cli.Context) error {
			server.ListenAndServe()
			return nil
		},
	}
}
