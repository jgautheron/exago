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
	ScoreValue  float64
	WeightValue float64
	Msg         string
}

// Calculate computes the criteria evaluation score
func (c *Evaluator) Calculate(p map[string]interface{}) {
}

// Score returns the criteria evaluation score
func (c *Evaluator) Score() float64 {
	return c.ScoreValue
}

// Weight returns the criteria evaluationvi tou weight
func (c *Evaluator) Weight() float64 {
	return c.WeightValue
}

// Message returns the criteria evaluation message
func (c *Evaluator) Message() string {
	return c.Msg
}
