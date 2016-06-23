package main

import (
	"regexp"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/indexer"
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
				Usage: "index the passed repositories",
				Action: func(c *cli.Context) error {
					items := []string{}
					for _, item := range c.Args() {
						items = append(items, item)
					}

					idx := indexer.New(items)
					idx.Start()
					return nil
				},
			},
			{
				Name:  "godoc",
				Usage: "parse and index the entire Godoc.org index",
				Action: func(c *cli.Context) error {
					parseAndIndexGodoc()
					return nil
				},
			},
		},
	}
}

func parseAndIndexGodoc() {
	const GodocIndexURL = "https://godoc.org/-/index"

	doc, err := goquery.NewDocument(GodocIndexURL)
	if err != nil {
		log.Fatal(err)
	}

	r, _ := regexp.Compile(`^github.com/([\w\d\-]+)/([\w\d\-]+)`)

	out := map[string]bool{}
	doc.Find("td a").Each(func(i int, s *goquery.Selection) {
		matches := r.FindStringSubmatch(s.Contents().Text())
		if len(matches) == 0 {
			return
		}
		out[matches[0]] = true
	})

	log.Infof("Found %d unique GitHub repositories in the Godoc index", len(out))

	sl := []string{}
	for item, _ := range out {
		sl = append(sl, item)
	}

	idx := indexer.New(sl)
	idx.Start()
}
