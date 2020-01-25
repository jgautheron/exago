package cov

import (
	"bufio"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/cover"
	"simonwaldherr.de/go/golibs/xmath"
)

// Report contains information about tested packages, functions and statements
type Report struct {
	// Packages holds all tested packages
	Packages []*Package `json:"packages"`
	// Coverage
	Coverage float64 `json:"coverage"`
}

func (r *Report) parseProfile(profiles []*cover.Profile) error {
	conv := converter{
		packages: make(map[string]*Package),
	}
	for _, p := range profiles {
		if err := conv.convertProfile(p); err != nil {
			return err
		}
	}
	for _, pkg := range conv.packages {
		r.addPackage(pkg)
	}

	return nil
}

// collectPackages collects ALL packages
func (r *Report) collectPackages() error {
	set := token.NewFileSet()
	dirs, err := packageList("Dir")
	if err != nil {
		return err
	}

	var errs []string
	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(set, dir, nil, 0)
		if err != nil {
			err := fmt.Sprintf("Directory %s returned error: `%s`", dir, err.Error())
			logrus.Error(err)
			errs = append(errs, err)
		}
		for _, pkg := range pkgs {
			// Ignore test package
			if strings.HasSuffix(pkg.Name, "_test") {
				logrus.Debugf("Ignoring test package `%s`", pkg.Name)
				continue
			}
			// Craft package path
			path := strings.Replace(dir, os.Getenv("GOPATH")+"/src/", "", 1)

			logrus.Debugf("path %v / package %v", path, pkg.Name)
			p := &Package{
				Name: pkg.Name,
				Path: path,
			}
			// Count LOCs for each file in the package
			for fn := range pkg.Files {
				if strings.HasSuffix(fn, "_test.go") {
					continue
				}
				p.LOC += countLOC(fn)
			}
			r.addPackage(p)
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

// countLOC counts lines of code, pull LOC, Comments, assertions
func countLOC(path string) int {
	var (
		loc            int
		inBlockComment bool
	)

	f, err := os.Open(path)
	if err != nil {
		logrus.Error(err)
		return loc
	}
	defer f.Close()

	buff := bufio.NewReader(f)

	for {
		line, isPrefix, err := buff.ReadLine()
		if err == io.EOF {
			return loc
		}
		// Incomplete line (don't count)
		if isPrefix == true {
			continue
		}
		// Empty line (don't count)
		if len(line) == 0 {
			continue
		}
		// Comment (don't count)
		if strings.Index(strings.TrimSpace(string(line)), "//") == 0 {
			continue
		}

		blockCommentStartPos := strings.Index(strings.TrimSpace(string(line)), "/*")
		blockCommentEndPos := strings.LastIndex(strings.TrimSpace(string(line)), "*/")

		if blockCommentStartPos > -1 {
			// block was started and not terminated
			if blockCommentEndPos == -1 || blockCommentStartPos > blockCommentEndPos {
				inBlockComment = true
			}
		}
		if blockCommentEndPos > -1 {
			// end of block is found and no new block was started
			if blockCommentStartPos == -1 || blockCommentEndPos > blockCommentStartPos {
				inBlockComment = false
			}
		}

		if inBlockComment {
			continue
		}

		loc++
	}
}

// AddPackage adds a package coverage information
func (r *Report) addPackage(p *Package) {
	i := sort.Search(len(r.Packages), func(i int) bool {
		return (r.Packages)[i].Name >= p.Name
	})
	if i < len(r.Packages) && (r.Packages)[i].Name == p.Name {
		(r.Packages)[i].Accumulate(p)
	} else {
		head := (r.Packages)[:i]
		tail := append([]*Package{p}, (r.Packages)[i:]...)
		r.Packages = append(head, tail...)
	}
}

// computeGlobalCoverage compute the global coverage
// from all packages coverage, we use the weighted
// average for that calculating the weight proportionally
// based on the LOCs
func (r *Report) computeGlobalCoverage() {
	var sums, weights []float64

	// Make the SUM of all LOCs
	var gloc float64
	for _, pkg := range r.Packages {
		gloc += float64(pkg.LOC)
	}

	for _, pkg := range r.Packages {
		// Weight will be the ratio of {PACKAGE_LOC}/{GLOBAL_PACKAGE_LOC}
		w := float64(pkg.LOC) / gloc
		weights = append(weights, w)
		sums = append(sums, w*pkg.Coverage)
	}

	r.Coverage = xmath.Sum(sums) / xmath.Sum(weights)
}

// packageList returns a list of Go-like files or directories from PWD,
func packageList(arg string) ([]string, error) {
	cmd, err := exec.Command("sh", "-c", `go list -f '{{.`+arg+`}}' ./... | grep -v vendor | grep -v Godeps`).CombinedOutput()
	if err != nil {
		return nil, err
	}

	pl := strings.Split(strings.TrimSpace(string(cmd)), "\n")

	return pl, nil
}
