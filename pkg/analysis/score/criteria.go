// Package score defines interface and generic receiver to be implemented by each
// criteria evaluator
package score

import (
	exago "github.com/jgautheron/exago/pkg"
)

// CriteriaEvaluator is the interface that must be implemented by a criteria
// evaluator.
type CriteriaEvaluator interface {
	Calculate(exago.Data) *exago.EvaluatorResponse
	Name() string
	Setup()
}

// Evaluator is a type that implements CriteriaEvaluator by allowing nil
// values but otherwise delegating to another ValueConverter.
type Evaluator struct {
	name string
	url  string
	desc string
}

// NewResponse creates an EvaluatorResponse instance
func (c *Evaluator) NewResponse(score float64, weight float64, msg string, details []*exago.EvaluatorResponse) *exago.EvaluatorResponse {
	return &exago.EvaluatorResponse{c.Name(), score, weight, c.desc, msg, c.url, details}
}

// Calculate computes the criteria evaluation score
func (c *Evaluator) Calculate(d exago.Data) *exago.EvaluatorResponse {
	return nil
}

// Setup is called before Calculate
func (c *Evaluator) Setup() {

}

// Name returns the criteria name
func (c *Evaluator) Name() string {
	return c.name
}

// Desc returns the criteria description
func (c *Evaluator) Desc() string {
	return c.desc
}

// URL returns the criteria URL
func (c *Evaluator) URL() string {
	return c.url
}
