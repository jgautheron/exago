package checklist

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileType int

const (
	TypeDir FileType = iota
	TypeFile
	TypeBoth
)

// FilesExistAny checks if the given file(s) exists in the root folder.
func FilesExistAny(path string, tp FileType, files ...string) bool {
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}

	matchesWith := func(f os.FileInfo, files []string) bool {
		for _, file := range files {
			if strings.Index(strings.ToLower(f.Name()), file) != -1 {
				return true
			}
		}
		return false
	}

	for _, f := range dirFiles {
		switch tp {
		case TypeDir:
			if f.IsDir() && matchesWith(f, files) {
				return true
			}
		case TypeFile:
			if !f.IsDir() && matchesWith(f, files) {
				return true
			}
		case TypeBoth:
			if !f.IsDir() && matchesWith(f, files) {
				return true
			}
		}
	}

	return false
}

// FindOccurrencesInTree tries to match the regular expression in files matching the file pattern.
// It returns the number of matchings.
func FindOccurrencesInTree(path, regex, filePattern string) int {
	matches := 0

	err := filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only interested in files
		if f.IsDir() {
			return nil
		}

		if match, err := filepath.Match(filePattern, f.Name()); !match || err != nil {
			return nil
		}

		file, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		r, _ := regexp.Compile(regex)
		match := r.FindStringSubmatch(string(file))

		if len(match) > 0 {
			matches++
		}
		return nil
	})

	if err != nil {
		return 0
	}
	return matches
}
