package model

const LintMessagesName = "lintmessages"

// LintMessages stores messages returned by Go linters
type LintMessages map[string]map[string][]map[string]interface{}
