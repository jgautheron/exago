// Package indexer enables mass processing of repositories.
package indexer

import "github.com/hotolab/exago-svc/queue"

// Start runs the indexer.
func Start(repos []string) {
	for _, repo := range repos {
		queue.GetInstance().PushAsync(repo, 10)
	}
	// queue.WaitUntilEmpty()
}
