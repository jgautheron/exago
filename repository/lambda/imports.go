package lambda

import (
	"encoding/json"

	"github.com/exago/svc/repository/model"
)

var importsCmd = &cmd{
	name:      new(model.Imports).Name(),
	unMarshal: unMarshalImports,
}

func GetImports(repository string) (model.RepositoryData, error) {
	importsCmd.ctxt = context{
		Repository: repository,
	}
	return importsCmd.Data()
}

func unMarshalImports(l *cmd, b []byte) (model.RepositoryData, error) {
	var imp model.Imports
	err := json.Unmarshal(b, &imp)
	return imp, err
}
