package main

import (
	"errors"

	. "github.com/hotolab/exago-svc"
	"github.com/hotolab/exago-svc/github"
	"github.com/hotolab/exago-svc/godoc"
	"github.com/hotolab/exago-svc/leveldb"
	"github.com/hotolab/exago-svc/pool"
	"github.com/hotolab/exago-svc/pool/job"
	"github.com/hotolab/exago-svc/repository/loader"
	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/repository/processor"
	"github.com/urfave/cli"
)

var db model.Database

// IndexCommand saves the godoc index in DB.
func IndexCommand() cli.Command {
	return cli.Command{
		Name:  "index",
		Usage: "Index repositories in the database",
		Subcommands: []cli.Command{
			{
				Name:  "repos",
				Usage: "Index the passed repositories",
				Action: func(c *cli.Context) error {
					items := []string{}
					for _, item := range c.Args() {
						items = append(items, item)
					}
					pl, err := initPool()
					if err != nil {
						return err
					}
					indexRepos(pl, items)
					return nil
				},
			},
			{
				Name:  "godoc",
				Usage: "Parse and index the entire Godoc.org index",
				Action: func(c *cli.Context) error {
					return indexGodoc()
				},
			},
		},
	}
}

func initPool() (pl model.Pool, err error) {
	// Initialise the lambda connection
	job.New()

	db, err = leveldb.New()
	if err != nil {
		return nil, err
	}

	rh, err := github.New()
	if err != nil {
		return nil, err
	}

	rl := loader.New(
		WithDatabase(db),
		WithRepositoryHost(rh),
	)

	po := processor.New(
		WithRepositoryLoader(rl),
	)

	return pool.New(
		WithProcessor(po.ProcessRepository),
	)
}

func indexGodoc() error {
	pl, err := initPool()
	if err != nil {
		return err
	}

	repos, err := godoc.New(
		WithDatabase(db),
	).GetIndex()
	if err != nil {
		return errors.New("Got error while trying to load the repos, did you index before godoc?")
	}

	repos = repos[:6]
	indexRepos(pl, repos)
	return nil
}

func indexRepos(pl model.Pool, repos []string) {
	for _, repo := range repos {
		pl.PushAsync(repo)
	}
	pl.WaitUntilEmpty()
}
