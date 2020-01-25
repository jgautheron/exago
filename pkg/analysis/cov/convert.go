package cov

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"
)

type converter struct {
	packages map[string]*Package
}

type extent struct {
	startOffset int
	startLine   int
	startCol    int
	endOffset   int
	endLine     int
	endCol      int
	coverage    float64
}

// convertProfile converts a Go coverage profile into an intelligent
// structure containing the percent of coverage, etc...
func (c *converter) convertProfile(p *cover.Profile) error {
	name, file, pkgpath, abspath, err := c.findFile(p.FileName)
	if err != nil {
		return err
	}
	pkg := c.packages[name]
	if pkg == nil {
		pkg = &Package{Name: name, Path: pkgpath}
		c.packages[name] = pkg
	}
	// Find function and statement extents; create corresponding
	// cov.Functions and cov.Statements, and keep a separate
	// slice of gocov.Statements so we can match them with profile
	// blocks.
	extents, err := c.findFuncs(file)
	if err != nil {
		return err
	}
	var stmts []statement
	for _, fe := range extents {
		f := &Function{
			Name:  fe.name,
			File:  abspath,
			Start: fe.startLine,
			End:   fe.endLine,
		}
		for _, stmt := range fe.stmts {
			s := statement{
				Statement:  &Statement{Start: stmt.startLine, End: stmt.endLine},
				StmtExtent: stmt,
			}
			f.Statements = append(f.Statements, s.Statement)
			stmts = append(stmts, s)
		}

		pkg.Functions = append(pkg.Functions, f)
	}

	// For each profile block in the file, find the statement(s) it
	// covers and increment the Reached field(s).
	blocks := p.Blocks
	for len(stmts) > 0 {
		s := stmts[0]
		for i, b := range blocks {
			if b.StartLine > s.endLine || (b.StartLine == s.endLine && b.StartCol >= s.endCol) {
				// Past the end of the statement
				stmts = stmts[1:]
				blocks = blocks[i:]
				break
			}
			if b.EndLine < s.startLine || (b.EndLine == s.startLine && b.EndCol <= s.startCol) {
				// Before the beginning of the statement
				continue
			}

			s.Reached += int64(b.Count)

			stmts = stmts[1:]
			break
		}
	}

	// Loop on each statement and determine coverage and TLOC by function
	var totalStmts int
	var totalReached int64
	for _, fn := range pkg.Functions {
		var reached int64
		totalStmts += len(fn.Statements)
		for _, stmt := range fn.Statements {
			if stmt.Reached > 0 {
				reached++
			}
		}

		totalReached += reached
		fn.Coverage = 100.0 * float64(reached) / float64(len(fn.Statements))
	}

	pkg.Coverage = 100.0 * float64(totalReached) / float64(totalStmts)

	return nil
}

// findFile finds the location of the named file in GOROOT, GOPATH etc.
func (c *converter) findFile(file string) (pkgname string, filename string, pkgpath string, abspath string, err error) {
	dir, file := filepath.Split(file)
	if dir != "" {
		dir = dir[:len(dir)-1] // drop trailing '/'
	}
	pkg, err := build.Import(dir, ".", build.IgnoreVendor)
	if err != nil {
		return "", "", "", "", fmt.Errorf("can't find %q: %v", file, err)
	}

	dir = strings.Replace(filepath.Join(pkg.Dir, file), pkg.SrcRoot, "$GOPATH", 1)
	return pkg.Name, filepath.Join(pkg.Dir, file), strings.Replace(pkg.ImportPath, pkg.SrcRoot, "", 1), dir, nil
}

// findFuncs parses the file and returns a slice of FuncExtent descriptors.
func (c *converter) findFuncs(name string) ([]*FuncExtent, error) {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}
	visitor := &FuncVisitor{fset: fset}
	ast.Walk(visitor, parsedFile)

	return visitor.funcs, nil
}
