package main

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	. "github.com/hotolab/exago-svc"
	"github.com/hotolab/exago-svc/github"
	"github.com/hotolab/exago-svc/gosearch"
	"github.com/hotolab/exago-svc/leveldb"
	"github.com/hotolab/exago-svc/pool"
	"github.com/hotolab/exago-svc/pool/job"
	"github.com/hotolab/exago-svc/repository/loader"
	"github.com/hotolab/exago-svc/repository/model"
	"github.com/hotolab/exago-svc/repository/processor"
	"github.com/urfave/cli"
)

var (
	db model.Database
	rl model.RepositoryLoader
)

// IndexCommand saves the godoc index in DB.
func IndexCommand() cli.Command {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "goversion",
			Value: "1.7.4",
			Usage: "Go version to be used for processing repositories",
		},
	}
	return cli.Command{
		Name:  "index",
		Usage: "Index repositories in the database",
		Subcommands: []cli.Command{
			{
				Name:  "repos",
				Usage: "Index the passed repositories",
				Flags: flags,
				Action: func(c *cli.Context) error {
					items := []string{}
					for _, item := range c.Args() {
						items = append(items, item)
					}
					pl, err := initPool()
					if err != nil {
						return err
					}
					indexRepos(c, pl, items)
					return nil
				},
			},
			{
				Name:  "gosearch",
				Usage: "Process the entire Gosearch index",
				Flags: flags,
				Action: func(c *cli.Context) error {
					return indexGosearch(c)
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

	rl = loader.New(
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

func indexGosearch(c *cli.Context) error {
	pl, err := initPool()
	if err != nil {
		return err
	}

	repos, err := gosearch.New(
		WithDatabase(db),
	).GetIndex()
	if err != nil {
		return errors.New("Got error while trying to load the repos, did you index before godoc?")
	}

	indexRepos(c, pl, repos)
	return nil
}

func indexRepos(c *cli.Context, pl model.Pool, repos []string) {
	goversion := c.String("goversion")
	for _, repo := range repos {
		data, err := rl.IsValid(repo)
		if err != nil {
			log.WithFields(log.Fields{
				"repository": repo,
				"error":      err.Error(),
			}).Error("Skipped indexing")
			continue
		}
		defaultBranch := data["default_branch"].(string)
		if !rl.IsCached(repo, defaultBranch, goversion) {
			pl.PushAsync(repo, defaultBranch, goversion)
		}
	}
	pl.WaitUntilEmpty()
}
