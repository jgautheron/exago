package main

import (
	. "github.com/jgautheron/exago"
	"github.com/jgautheron/exago/gosearch"
	"github.com/jgautheron/exago/leveldb"
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
