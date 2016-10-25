package lambda

import (
	"encoding/json"

	"github.com/hotolab/exago-svc/repository/model"
)

func (l Runner) FetchCodeStats() (model.CodeStats, error) {
	codeStatsCmd := &cmd{
		name:      model.CodeStatsName,
		unMarshal: unMarshalCodeStats,
	}
	codeStatsCmd.ctxt = context{
		Repository: l.Repository,
		Cleanup:    l.ShouldCleanup,
	}
	d, err := codeStatsCmd.Data()
	if err != nil {
		return nil, err
	}
	return d.(model.CodeStats), nil
}

func unMarshalCodeStats(l *cmd, b []byte) (interface{}, error) {
	var cs model.CodeStats
	err := json.Unmarshal(b, &cs)
	return cs, err
}
