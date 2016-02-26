package datafetcher

import (
	"encoding/json"

	"github.com/exago/svc/repository"
)

var importsCmd = &lambdaCmd{
	name:      "imports",
	unMarshal: unMarshalImports,
}

func GetImports(repository string) (interface{}, error) {
	importsCmd.ctxt = lambdaContext{
		Repository: repository,
	}
	return importsCmd.Data()
}

func unMarshalImports(l *lambdaCmd, b []byte) (interface{}, error) {
	var imp repository.Imports
	err := json.Unmarshal(b, &imp)
	return imp, err
}
