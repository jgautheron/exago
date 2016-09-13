package model

const CodeStatsName = "codestats"

// CodeStats stores infos about code such as ratio of LOC vs CLOC etc..
type CodeStats map[string]int
