package cov

import (
	"os"

	"golang.org/x/tools/cover"
)

// ConvertRepository converts a given repository to a Report struct
func ConvertRepository(repo string) (*Report, error) {
	r := &Report{}
	err := r.collectPackages()
	if err != nil {
		return nil, err
	}

	p, err := createProfile()
	if err != nil {
		return nil, err
	}
	defer os.Remove(p.Name())

	profiles, err := cover.ParseProfiles(p.Name())
	if err != nil {
		return nil, err
	}

	if err = r.parseProfile(profiles); err != nil {
		return nil, err
	}

	r.computeGlobalCoverage()

	return r, nil
}
