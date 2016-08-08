package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/exago-svc/godoc"
	"github.com/urfave/cli"
)

// GodocCommand handles all godoc-related actions.
func GodocCommand() cli.Command {
	return cli.Command{
		Name:  "godoc",
		Usage: "Godoc related actions",
		Subcommands: []cli.Command{
			{
				Name:  "save",
				Usage: "Save the Godoc index in database",
				Action: func(c *cli.Context) error {
					if err := godoc.New().SaveIndex(); err != nil {
						return err
					}
					log.Info("Successfully persisted in DB the Godoc index")
					return nil
				},
			},
		},
	}
}
