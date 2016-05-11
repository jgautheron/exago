package repository

import "fmt"

const (
	thirdParties scoreType = "thirdParties"
	ratioLocCloc           = "ratioLocCloc"
	checklist              = "checklist"
	testPass               = "testPass"
	testCoverage           = "testCoverage"
	testDuration           = "testDuration"
)

type scoreType string

type Score struct {
	Value   int      `json:"value"`
	Rank    Rank     `json:"rank"`
	Details []string `json:"details"`
}

func (sc *Score) Increase(score int, category scoreType, comment string) {
	sc.addInfo(fmt.Sprintf("+%d: %s %s", score, category, comment))
	sc.Value += score
}

func (sc *Score) Decrease(score int, category scoreType, comment string) {
	sc.addInfo(fmt.Sprintf("-%d: %s %s", score, category, comment))
	sc.Value -= score
}

func (sc *Score) stampRank() {
	switch true {
	case sc.Value >= 80:
		sc.Rank = A
	case sc.Value >= 60:
		sc.Rank = B
	case sc.Value >= 40:
		sc.Rank = C
	case sc.Value >= 20:
		sc.Rank = D
	case sc.Value >= 0:
		sc.Rank = E
	default:
		sc.Rank = F
	}
}

func (sc *Score) addInfo(str string) {
	sc.Details = append(sc.Details, str)
}

// TODO: Split the score calculation per Struct (ex. Imports{}.GetScore())
func (r *Repository) calcScore() Score {
	// More third parties means bigger potential for instability, larger attack surface
	tp := len(r.Imports)
	switch true {
	case tp == 0:
		r.Score.Increase(20, thirdParties, "0")
	case tp < 4:
		r.Score.Increase(15, thirdParties, "< 4")
	case tp < 6:
		r.Score.Increase(10, thirdParties, "< 6")
	case tp < 8:
		r.Score.Increase(5, thirdParties, "< 8")
	}

	// Code doesn't always speak for itself
	ra := float64(r.CodeStats["LOC"] / r.CodeStats["NCLOC"])
	switch true {
	case ra > 1.4:
		r.Score.Increase(20, ratioLocCloc, "> 1.4")
	case ra > 1.3:
		r.Score.Increase(15, ratioLocCloc, "> 1.3")
	case ra > 1.2:
		r.Score.Increase(10, ratioLocCloc, "> 1.2")
	case ra > 1.1:
		r.Score.Increase(8, ratioLocCloc, "> 1.1")
	}

	// Checklist
	for _, passed := range r.TestResults.Checklist.Passed {
		switch passed.Name {
		case "projectBuilds":
			r.Score.Increase(10, checklist, "projectBuilds")
		case "isFormatted":
			r.Score.Increase(10, checklist, "isFormatted")
		case "hasReadme":
			r.Score.Increase(10, checklist, "hasReadme")
		case "isDirMatch":
			r.Score.Increase(10, checklist, "isDirMatch")
		case "isLinted":
			r.Score.Increase(10, checklist, "isLinted")
		case "hasBenches":
			r.Score.Increase(10, checklist, "hasBenches")
		}
	}
	for _, failed := range r.TestResults.Checklist.Failed {
		switch failed.Name {
		case "isFormatted":
			r.Score.Decrease(25, checklist, "isFormatted")
		case "isLinted":
			r.Score.Decrease(10, checklist, "isLinted")
		case "hasReadme":
			r.Score.Decrease(20, checklist, "hasReadme")
		case "isDirMatch":
			r.Score.Decrease(10, checklist, "isDirMatch")
		}
	}

	// Initialise values from test results
	var cov, duration []float64
	success := true
	for _, pkg := range r.TestResults.Packages {
		cov = append(cov, pkg.Coverage)
		duration = append(duration, pkg.ExecutionTime)

		if !pkg.Success {
			success = false
		}
	}

	// Calculate mean values for both code coverage and execution time
	var covMean, durationMean float64 = 0, 0
	if len(cov) > 0 {
		for _, v := range cov {
			covMean += v
		}
		covMean /= float64(len(cov))
	}
	if len(duration) > 0 {
		for _, v := range duration {
			durationMean += v
		}
		durationMean /= float64(len(duration))
	}

	if !success {
		r.Score.Decrease(15, testPass, "")
	}

	// 100% is not necessarily an achievement
	switch true {
	case covMean > 80:
		r.Score.Increase(20, testCoverage, "> 80")
	case covMean > 60:
		r.Score.Increase(15, testCoverage, "> 60")
	case covMean > 40:
		r.Score.Increase(10, testCoverage, "> 40")
	case covMean == 0:
		r.Score.Decrease(20, testCoverage, "= 0")
	}

	// Fast test suites are important
	switch true {
	case durationMean > 10:
		r.Score.Decrease(15, testDuration, "> 10")
	case durationMean < 2:
		r.Score.Increase(10, testDuration, "< 2")
	}

	// Stamp the rank
	r.Score.stampRank()

	return r.Score
}
