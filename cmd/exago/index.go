package main

import (
	"time"

	"github.com/hotolab/exago-svc/github"
	"github.com/hotolab/exago-svc/pool"
	"github.com/hotolab/exago-svc/repository/processor"
	"github.com/hotolab/exago-svc/taskrunner/lambda"
	"github.com/pkg/profile"
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
	// lambda.GetInstance()

	pool.TestCustomWorkers()

	// repos, err := godoc.New().GetIndex()
	// if err != nil {
	// 	return errors.New("Got error while trying to load the repos, did you index before godoc?")
	// }
	// indexRepos(repos)
	return nil
}

func indexRepos(repos []string) {
	// github.GetInstance()
	// lambda.GetInstance()
	// q := queue.GetInstance()
	list := repos[:30]
	defer profile.Start(profile.BlockProfile).Stop()
	for _, repo := range list {
		go process(repo)
	}
	time.Sleep(6 * time.Minute)
	// q.WaitUntilEmpty()
}

func process(repo string) {
	processor.ProcessRepository(repo, "", lambda.Runner{Repository: repo})
}
