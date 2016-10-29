package model

import "time"

type Data struct {
	CodeStats     CodeStats         `json:"codestats"`
	ProjectRunner ProjectRunner     `json:"projectrunner"`
	LintMessages  LintMessages      `json:"lintmessages"`
	Metadata      Metadata          `json:"metadata"`
	Score         Score             `json:"score"`
	Name          string            `json:"name"`
	Branch        string            `json:"branch"`
	ExecutionTime string            `json:"execution_time"`
	LastUpdate    time.Time         `json:"last_update"`
	Errors        map[string]string `json:"errors,omitempty"`
}
