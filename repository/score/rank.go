package score

const (
	A string = "A"
	B        = "B"
	C        = "C"
	D        = "D"
	E        = "E"
	F        = "F"
)

func Rank(value float64) string {
	switch true {
	case value >= 80:
		return A
	case value >= 60:
		return B
	case value >= 40:
		return C
	case value >= 20:
		return D
	case value >= 0:
		return E
	default:
		return F
	}
}
