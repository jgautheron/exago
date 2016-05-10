package lambda

import (
	"encoding/json"

	"github.com/exago/svc/repository/model"
)

var codeStatsCmd = &cmd{
	name:      new(model.CodeStats).Name(),
	unMarshal: unMarshalCodeStats,
}

func GetCodeStats(repository string) (model.RepositoryData, error) {
	codeStatsCmd.ctxt = context{
		Repository: repository,
	}
	return codeStatsCmd.Data()
}

func unMarshalCodeStats(l *cmd, b []byte) (model.RepositoryData, error) {
	var cs model.CodeStats
	err := json.Unmarshal(b, &cs)
	return cs, err
}
