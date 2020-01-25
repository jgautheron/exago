package main

import (
	"github.com/jgautheron/exago/pkg/analysis/task"
)

func main() {
	repo := "github.com/pkg/errors"
	m := task.NewManager(repo)

	//m.UseReference(c.String("ref"))

	// 1. run tasks
	// 2. import output to cov
	// 3. save in firestore

	res := m.ExecuteRunners()
	if res.Success {

	}
	//enc := json.NewEncoder(os.Stdout)
	//enc.Encode(out)
}
