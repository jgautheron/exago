package main

import (
	"fmt"
	"time"

	"github.com/hotolab/exago-svc/godoc"
	"github.com/hotolab/exago-svc/queue"
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
		return fmt.Errorf("Got error while trying to load repos from GitHub: %v", err)
	}
	indexRepos(repos)
	return nil
}

func indexRepos(repos []string) {
	q := queue.GetInstance()
	list := repos[:5]
	for _, repo := range list {
		q.PushAsync(repo)
	}
	time.Sleep(1 * time.Second)
	q.WaitUntilEmpty()
}
