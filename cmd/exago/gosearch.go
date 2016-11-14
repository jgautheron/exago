package main

import (
	. "github.com/hotolab/exago-svc"
	"github.com/hotolab/exago-svc/gosearch"
	"github.com/hotolab/exago-svc/leveldb"
	"github.com/urfave/cli"
)

// GosearchCommand handles all gosearch-related actions.
func GosearchCommand() cli.Command {
	return cli.Command{
		Name:  "gosearch",
		Usage: "GoSearch related actions",
		Subcommands: []cli.Command{
			{
				Name:  "save",
				Usage: "Save the GoSearch index in database",
				Action: func(c *cli.Context) error {
					db, err := leveldb.New()
					if err != nil {
						return err
					}
					return gosearch.New(
						WithDatabase(db),
					).LoadIndex()
				},
			},
		},
	}
}
