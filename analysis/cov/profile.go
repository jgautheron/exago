package cov

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	coverMode = "count"
)

// processPackage executes go test command with coverage and outputs
// errors and output into channels so they are combined later in a single
// file and passed to cov for getting the expected JSON output
func processPackage(rel string) (string, error) {
	// Create temporary file to output the file coverage
	// this file is trashed after processing
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	logrus.Debugf("go test -covermode=%s -coverprofile=%s %s", coverMode, tmp.Name(), rel)
	_, err = exec.Command("go", "test", "-covermode="+coverMode, "-coverprofile="+tmp.Name(), rel).CombinedOutput()
	if err != nil {
		return "", nil
	}

	// Get file contents
	b, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// lookupTestFiles crawls the filesystem from the repository path
// and finds test files using glob, if a package doesn't have tests
// it is automatically skipped.
func createProfile() (*os.File, error) {
	pkgs, err := packageList("ImportPath")
	if err != nil {
		return nil, err
	}

	file, err := ioutil.TempFile("", "hotolab-coverage")
	if err != nil {
		return nil, err
	}

	// Bufferize channel
	tasks := make(chan string, 64)
	var (
		wg      sync.WaitGroup
		errBuff bytes.Buffer
		outBuff bytes.Buffer
	)

	// Create as much threads as we have CPUs
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			for pkg := range tasks {
				res, err := processPackage(pkg)
				if err != nil {
					errBuff.WriteString(err.Error())
					return
				}
				outBuff.WriteString(res)
			}
			wg.Done()
		}()
	}

	for _, pkg := range pkgs {
		tasks <- pkg
	}

	// Close worker channel
	close(tasks)

	// Wait for the workers to finish
	wg.Wait()

	// Get errors (if any) and convert them to a runner error
	errs := errBuff.String()
	if errs != "" {
		return nil, errors.New(errs)
	}

	// Get content of the buffer and write it
	// to the temp file attached to the runner
	out := outBuff.String()
	out = regexp.MustCompile("mode: [a-z]+\n").ReplaceAllString(out, "")
	out = "mode: " + coverMode + "\n" + out

	err = ioutil.WriteFile(file.Name(), []byte(out), 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}
