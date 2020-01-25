package cov

import "fmt"

// Package describes a package inner characteristics
type Package struct {
	// Name is the package name
	Name string `json:"name"`
	// Path is the canonical path of the package.
	Path string `json:"path"`
	// Coverage
	Coverage float64 `json:"coverage"`
	// LOC contains the number of lines of code for a given package
	LOC int `json:"loc"`
	// Functions is a list of functions registered with this package.
	Functions []*Function `json:"-"`
}

// Accumulate will accumulate the coverage information from the provided
// Package into this Package.
func (p *Package) Accumulate(p2 *Package) error {
	if p.Name != p2.Name {
		return fmt.Errorf("Names do not match: %q != %q", p.Name, p2.Name)
	}
	if p.Coverage != p2.Coverage {
		p.Coverage = p2.Coverage
	}
	if p.Path != p2.Path {
		p.Path = p2.Path
	}
	if len(p.Functions) != len(p2.Functions) {
		return fmt.Errorf("Function counts do not match: %d != %d", len(p.Functions), len(p2.Functions))
	}
	for i, f := range p.Functions {
		err := f.Accumulate(p2.Functions[i])
		if err != nil {
			return err
		}
	}

	return nil
}
