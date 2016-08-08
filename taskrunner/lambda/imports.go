package lambda

import (
	"encoding/json"
	"regexp"

	"github.com/hotolab/exago-svc/repository/model"
)

var importsCmd = &cmd{
	name:      model.ImportsName,
	unMarshal: unMarshalImports,
}

func (l Runner) FetchImports() (model.Imports, error) {
	importsCmd.ctxt = context{
		Repository: l.Repository,
		Cleanup:    l.ShouldCleanup,
	}
	d, err := importsCmd.Data()
	if err != nil {
		return nil, err
	}
	data := d.(model.Imports)

	// Dedupe third party packages
	// One repository corresponds to one third party
	imports, filtered := []string{}, map[string]int{}
	reg, _ := regexp.Compile(`^github\.com/([\w\d\-]+)/([\w\d\-]+)`)
	for _, im := range data {
		m := reg.FindStringSubmatch(im)
		if len(m) > 0 {
			filtered[m[0]] = 1
		} else {
			filtered[im] = 1
		}
	}
	for im := range filtered {
		imports = append(imports, im)
	}

	return imports, nil
}

func unMarshalImports(l *cmd, b []byte) (interface{}, error) {
	var imp model.Imports
	err := json.Unmarshal(b, &imp)
	return imp, err
}
