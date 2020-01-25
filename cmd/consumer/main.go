package main

import (
	"encoding/json"
	"fmt"

	exago "github.com/jgautheron/exago/pkg"

	"github.com/jgautheron/exago/pkg/analysis/task"
)

func main() {
	repo := "github.com/pkg/errors"
	m := task.NewManager(repo)

	//m.UseReference(c.String("ref"))

	res := m.ExecuteRunners()
	if res.Success {
		fmt.Printf("%#v", res.Runners)
	}

	out, err := json.Marshal(res.Runners)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))

	var foo exago.Results
	err = json.Unmarshal(out, &foo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", foo)
}
