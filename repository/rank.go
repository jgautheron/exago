package repository

const (
	A Rank = "A"
	B      = "B"
	C      = "C"
	D      = "D"
	E      = "E"
	F      = "F"
)

type Rank string

func (r *Repository) Rank() Rank {
	return r.Score.Rank
}
