package score_test

import (
	"testing"

	"github.com/jgautheron/exago/analysis/score"
)

func TestRank(t *testing.T) {
	var tests = []struct {
		score    float64
		expected string
	}{
		{99, "A+"},
		{95, "A"},
		{91, "A-"},
		{89, "B+"},
		{85, "B"},
		{80, "B-"},
		{79, "C+"},
		{75, "C"},
		{70, "C-"},
		{69, "D+"},
		{66, "D"},
		{60, "D-"},
		{59, "E+"},
		{54, "E"},
		{50, "E-"},
		{49, "F+"},
		{45, "F"},
		{40, "F-"},
		{30, "F-"},
		{20, "F-"},
		{10, "F-"},
		{0, "F-"},
	}

	for _, tt := range tests {
		rank := score.Rank(tt.score)
		if rank != tt.expected {
			t.Error("Wrong rank")
		}
	}
}
