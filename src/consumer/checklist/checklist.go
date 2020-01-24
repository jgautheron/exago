package checklist

import (
	"os"
	"strings"
	"sync"
)

type CheckList struct {
	checkList    []CheckItem
	sourcePath   string
	sourceGoPath string
}

func New(sourcePath string) *CheckList {
	sourceGoPath := strings.Replace(sourcePath, os.Getenv("GOPATH")+"/src/", "", 1)

	checkList := []CheckItem{
		{
			Name: "isFormatted",
			Desc: "gofmt Correctness: Is the code formatted correctly?",
			fn:   isFormatted,
		},
		{
			Name: "isLinted",
			Desc: "golint Correctness: Is the linter satisfied?",
			fn:   isLinted,
		},
		{
			Name: "isVetted",
			Desc: "go tool vet Correctness: Is the Go vet satisfied?",
			fn:   isVetted,
		},
		{
			Name: "hasLicense",
			Desc: "Licensed: Does the project have a license?",
			fn:   hasFiles(TypeFile, "license"),
		},
		{
			Name: "hasReadme",
			Desc: "README Presence: Does the project's include a documentation entrypoint?",
			fn:   hasFiles(TypeFile, "readme"),
		},
		{
			Name: "hasCI",
			Desc: "Is the project using a CI tool?",
			fn:   hasFiles(TypeFile, "circle.yml"),
		},
		{
			Name: "hasOldDep",
			Desc: "Is the project using outdated dependencies manager?",
			fn:   hasFiles(TypeBoth, "glide", "Godeps", "Gopkg"),
		},
		{
			Name: "hasGoMod",
			Desc: "Is the project using go.mod?",
			fn:   hasFiles(TypeFile, "go.mod"),
		},
		{
			Name: "hasContributing",
			Desc: "Contribution Process: Does the project document a contribution process?",
			fn:   hasFiles(TypeFile, "contribution", "contribute", "contributing"),
		},
		{
			Name: "hasChangelog",
			Desc: "Is the project maintaining a changelog?",
			fn:   hasFiles(TypeFile, "changelog"),
		},
		{
			Name: "hasBenches",
			Desc: "Benchmarks: In addition to tests, does the project have benchmarks?",
			fn:   hasOccurrence(`func\sBenchmark\w+\(`, "*_test.go"),
		},
		{
			Name: "hasMainPackage",
			Desc: "Does the project have a main package?",
			fn:   hasOccurrence(`"package(\s)main"`, "*.go"),
		},
		//{
		//	Name: "hasBlackboxTests",
		//	Desc: "Blackbox Tests: In addition to standard tests, does the project have blackbox tests?",
		//	fn:   hasOccurrence(`"testing\/quick"`, "*_test.go"),
		//},
	}

	return &CheckList{checkList, sourcePath, sourceGoPath}
}

// RunTasks is a wrapper for running all tasks from the list
func (c CheckList) RunTasks() (successful []string, failed []string) {
	var wg sync.WaitGroup

	wg.Add(len(c.checkList))
	for _, task := range c.checkList {
		go func(task CheckItem) {
			if ok := task.run(c.sourcePath, c.sourceGoPath); ok {
				successful = append(successful, task.Name)
			} else {
				failed = append(failed, task.Name)
			}
			wg.Done()
		}(task)
	}

	wg.Wait()
	return successful, failed
}
