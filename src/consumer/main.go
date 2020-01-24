package main

import (
	"encoding/json"
	"os"

	"github.com/jgautheron/exago/src/consumer/task"
)

func main() {
	repo := "github.com/pkg/errors"
	m := task.NewManager(repo)

	//m.UseReference(c.String("ref"))

	out := m.ExecuteRunners()
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(out)
}
