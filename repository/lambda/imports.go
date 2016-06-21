package lambda

import (
	"encoding/json"

	"github.com/exago/svc/repository/model"
)

var importsCmd = &cmd{
	name:      model.ImportsName,
	unMarshal: unMarshalImports,
}

func GetImports(repository string) (interface{}, error) {
	importsCmd.ctxt = context{
		Repository: repository,
	}
	return importsCmd.Data()
}

func unMarshalImports(l *cmd, b []byte) (interface{}, error) {
	var imp model.Imports
	err := json.Unmarshal(b, &imp)
	return imp, err
}
