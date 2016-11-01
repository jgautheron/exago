package main

import (
	. "github.com/hotolab/exago-svc"
	"github.com/hotolab/exago-svc/github"
	"github.com/hotolab/exago-svc/leveldb"
	"github.com/hotolab/exago-svc/pool"
	"github.com/hotolab/exago-svc/pool/job"
	"github.com/hotolab/exago-svc/repository"
	"github.com/hotolab/exago-svc/repository/processor"
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
			// Initialise the lambda connection
			job.New()

			db, err := leveldb.New()
			if err != nil {
				return err
			}

			rh, err := github.New()
			if err != nil {
				return err
			}

			rl := repository.NewLoader(
				WithDatabase(db),
				WithRepositoryHost(rh),
			)

			po := processor.New(
				WithRepositoryLoader(rl),
			)

			pl, err := pool.New(
				WithProcessor(po.ProcessRepository),
			)
			if err != nil {
				return err
			}

			sh, err := showcaser.New(
				WithDatabase(db),
				WithRepositoryLoader(rl),
			)
			if err != nil {
				return err
			}
			sh.StartRoutines()

			s := server.New(
				WithDatabase(db),
				WithRepositoryHost(rh),
				WithPool(pl),
				WithShowcaser(sh),
				WithRepositoryLoader(rl),
			)
			return s.ListenAndServe()
		},
	}
}
