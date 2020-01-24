package checklist

import (
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/lint"
)

func isFormatted() CheckItemParams {
	return func(sourcePath, sourceGoPath string) bool {
		errors := 0
		filepath.Walk(sourcePath, func(path string, f os.FileInfo, err error) error {
			if !strings.HasSuffix(filepath.Ext(path), ".go") {
				return nil
			}

			file, err := ioutil.ReadFile(path)
			if err != nil {
				return nil
			}

			fmtFile, _ := format.Source(file)

			if string(file) != string(fmtFile) {
				errors++
			}
			return nil
		})
		return errors == 0
	}
}

func isLinted() CheckItemParams {
	return func(sourcePath, sourceGoPath string) bool {
		errors := 0
		l := new(lint.Linter)

		filepath.Walk(sourcePath+"/...", func(path string, f os.FileInfo, err error) error {

			if !strings.HasSuffix(filepath.Ext(path), ".go") {
				return nil
			}

			file, err := ioutil.ReadFile(path)
			if err != nil {
				return nil
			}

			if lnt, _ := l.Lint(f.Name(), file); len(lnt) > 0 {
				if lnt[0].Confidence > 0.2 {
					errors++
					return nil
				}
			}
			return nil
		})

		return errors == 0
	}
}

func isVetted() CheckItemParams {
	return func(sourcePath, sourceGoPath string) bool {
		_, err := exec.Command("go", "vet", sourceGoPath).Output()
		return err == nil
	}
}

func hasFiles(tp FileType, files ...string) func() CheckItemParams {
	return func() CheckItemParams {
		return func(sourcePath, sourceGoPath string) bool {
			return FilesExistAny(sourcePath, tp, files...)
		}
	}
}

func hasOccurrence(regex, filePattern string) func() CheckItemParams {
	return func() CheckItemParams {
		return func(sourcePath, sourceGoPath string) bool {
			return FindOccurrencesInTree(sourcePath, regex, filePattern) > 0
		}
	}
}
