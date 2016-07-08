package model

import "time"

type Data struct {
	CodeStats     CodeStats        `json:"codestats"`
	Imports       Imports          `json:"imports"`
	TestResults   TestResults      `json:"testresults"`
	LintMessages  LintMessages     `json:"lintmessages"`
	Metadata      Metadata         `json:"metadata"`
	Score         Score            `json:"score"`
	ExecutionTime string           `json:"execution_time"`
	LastUpdate    time.Time        `json:"last_update"`
	Errors        map[string]error `json:"errors,omitempty"`
}
