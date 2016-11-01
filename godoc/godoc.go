// Package godoc contains the Godoc Index scraping logic.
package godoc

import (
	"encoding/json"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	exago "github.com/hotolab/exago-svc"
)

const (
	GodocIndexURL    = "https://godoc.org/-/index"
	GodocDatabaseKey = "godoc:index"
)

type Godoc struct {
	config exago.Config
}

func New(options ...exago.Option) *Godoc {
	var gd Godoc
	for _, option := range options {
		option.Apply(&gd.config)
	}
	return &gd
}

// GetIndex retrieves the Godoc Index from the database, meaning that SaveIndex
// must have been called before.
func (g *Godoc) GetIndex() (repos []string, err error) {
	b, err := g.config.DB.Get([]byte(GodocDatabaseKey))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &repos)
	return repos, err
}

// SaveIndex scrapes the Godoc Index, takes only GitHub repositories and ignores
// the notion of packages per repository.
// Then the index is persisted in database for later use (HTTP/ indexing).
func (g *Godoc) SaveIndex() error {
	doc, err := goquery.NewDocument(GodocIndexURL)
	if err != nil {
		return err
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

	b, err := json.Marshal(sl)
	if err != nil {
		return err
	}
	return g.config.DB.Put([]byte(GodocDatabaseKey), b)
}
