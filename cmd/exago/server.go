package main

import (
	. "github.com/jgautheron/exago"
	"github.com/jgautheron/exago/github"
	"github.com/jgautheron/exago/leveldb"
	"github.com/jgautheron/exago/pool"
	"github.com/jgautheron/exago/pool/job"
	"github.com/jgautheron/exago/repository/loader"
	"github.com/jgautheron/exago/repository/processor"
	"github.com/jgautheron/exago/server"
	"github.com/jgautheron/exago/showcaser"
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

			rl := loader.New(
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
