package score

import "math"

const rankBoundsRate = 3.75

// Rank computes a rank given a decimal score
// Quite simply, we divide the score by 10 and use floor(x)
// to get the greatest integer value less than or equal to x
//
// Afterwards we get the rank by a map lookup with the pow index
// if the index is not found we return the worst rank :P
// We add a plus (+) or minus (-) sign to the rank if the score is
// either in lower or upper bound of the range.
func Rank(value float64) string {
	rnks := map[int]string{10: "A", 9: "A", 8: "B", 7: "C", 6: "D", 5: "E", 4: "F"}
	pow := int(math.Floor(value / 10))

	// Return worst rank
	if _, ok := rnks[pow]; !ok {
		return rnks[4] + "-"
	}

	// Get rank
	rnk := rnks[pow]
	// Calculate the mid-range e.g. 8 => 85
	tsh := float64(pow)*10 + 5

	if value >= tsh+rankBoundsRate {
		rnk += "+"
	} else if value <= tsh-rankBoundsRate {
		rnk += "-"
	}

	return rnk
}
