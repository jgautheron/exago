package main

import (
	"errors"

	"github.com/hotolab/exago-svc/godoc"
	"github.com/hotolab/exago-svc/pool"
	"github.com/hotolab/exago-svc/pool/job"
	"github.com/urfave/cli"
)

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

					indexRepos(items)
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

func indexGodoc() error {
	repos, err := godoc.New().GetIndex()
	if err != nil {
		return errors.New("Got error while trying to load the repos, did you index before godoc?")
	}
	indexRepos(repos)
	return nil
}

func indexRepos(repos []string) {
	job.Init()
	p := pool.GetInstance()
	for _, repo := range repos {
		p.PushAsync(repo)
	}
	p.WaitUntilEmpty()
}
