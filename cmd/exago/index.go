package main

import (
	"errors"
	"log"
	"time"

	"github.com/hotolab/exago-svc/github"
	"github.com/hotolab/exago-svc/godoc"
	"github.com/hotolab/exago-svc/queue"
	"github.com/hotolab/exago-svc/taskrunner/lambda"
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
	github.GetInstance()
	lambda.GetInstance()

	repos, err := godoc.New().GetIndex()
	log.Println(repos)
	if err != nil {
		return errors.New("Got error while trying to load the repos, did you index before godoc?")
	}
	indexRepos(repos)
	return nil
}

func indexRepos(repos []string) {
	github.GetInstance()
	lambda.GetInstance()
	q := queue.GetInstance()
	list := repos[:5]
	for _, repo := range list {
		q.PushAsync(repo)
	}
	time.Sleep(1 * time.Second)
	q.WaitUntilEmpty()
}
