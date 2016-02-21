package datafetcher

import "encoding/json"

type imports []string

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
	var imp imports
	err := json.Unmarshal(b, &imp)
	return imp, err
}
