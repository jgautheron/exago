package model

type Data struct {
	CodeStats     CodeStats         `json:"codestats"`
	ProjectRunner ProjectRunner     `json:"projectrunner"`
	LintMessages  LintMessages      `json:"lintmessages"`
	Metadata      Metadata          `json:"metadata"`
	Score         Score             `json:"score"`
	Errors        map[string]string `json:"errors,omitempty"`
}
