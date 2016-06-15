package main

import (
	"regexp"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/exago/svc/repository"
)

const (
	GodocIndexURL = "https://godoc.org/-/index"
)

// IndexCommand saves the godoc index in DB.
func IndexCommand() cli.Command {
	return cli.Command{
		Name:  "index",
		Usage: "Save the Godoc index in DB",
		Action: func(ctx *cli.Context) {
			indexGodoc()
		},
	}
}

func indexGodoc() {
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

	idx := indexer{
		QueueConcurrent: 20,
	}
	idx.addItems(out)
	idx.index()
}

type indexer struct {
	QueueItems      []item
	QueueConcurrent int
}

func (idx *indexer) addItems(items map[string]bool) {
	slice := []item{}
	for repo, _ := range items {
		slice = append(slice, item{
			repo,
			make(chan bool, 1),
			make(chan bool, 1),
		})
	}
	idx.QueueItems = slice
}

func (idx *indexer) index() {
	for i, item := range idx.QueueItems {
		lgr := log.WithFields(log.Fields{
			"repository": item.name,
			"index":      i,
		})
		lgr.Infof("Processing item...")
		go item.process()
		select {
		case <-item.skip:
			lgr.Infof("Item skipped")
		case <-item.done:
			lgr.Infof("Item processed")
		}
	}
}

type item struct {
	name       string
	done, skip chan bool
}

func (i *item) process() {
	lgr := log.WithField("repository", i.name)

	rc := repository.NewChecker(i.name)
	if rc.Repository.IsCached() {
		lgr.Debugf("Already cached")
		return
	}

	go func() {
		for err := range rc.Errors {
			log.WithField("error", err.Error()).Warn("Got an error, aborting...")
			rc.Abort()
			i.skip <- true
		}
	}()

	// Wait until the data is ready
	rc.Run()

	<-rc.Done
	i.done <- true
}
