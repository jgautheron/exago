package main

import (
	"github.com/hotolab/exago-svc/pool/job"
	"github.com/hotolab/exago-svc/server"
	"github.com/hotolab/exago-svc/showcaser"
	"github.com/urfave/cli"
)

// ServerCommand starts the HTTP server.
func ServerCommand() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start the HTTP server",
		Action: func(c *cli.Context) error {
			// Initialise the showcaser data
			showcaser.GetInstance()
			// Initialise the lambda connection
			job.Init()

			server.ListenAndServe()
			return nil
		},
	}
}
