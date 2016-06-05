// Package indexer processes repositories to determine occurrences, top-k, popularity.
package indexer

import (
	"github.com/exago/svc/repository"
)

type IndexedData struct {
	Recent  []string
	Top     []string
	Popular []string
}

func ProcessRepository(repo *repository.Repository) {

}
