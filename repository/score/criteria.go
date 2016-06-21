// Package score defines interface and generic receiver to be implemented by each
// criteria evaluator
package score

// CriteriaEvaluator is the interface that must be implemented by a criteria
// evaluator.
type CriteriaEvaluator interface {
	Calculate(map[string]interface{})
	Weight() float64
	Score() float64
	Message() string
}

// Evaluator is a type that implements CriteriaEvaluator by allowing nil
// values but otherwise delegating to another ValueConverter.
type Evaluator struct {
	score  float64
	weight float64
	msg    string
}

// Calculate computes the criteria evaluation score
func (c *Evaluator) Calculate(p map[string]interface{}) {
}

// Score returns the criteria evaluation score
func (c *Evaluator) Score() float64 {
	return c.score
}

// Weight returns the criteria evaluation weight
func (c *Evaluator) Weight() float64 {
	return c.weight
}

// Message returns the criteria evaluation message
func (c *Evaluator) Message() string {
	return c.msg
}
