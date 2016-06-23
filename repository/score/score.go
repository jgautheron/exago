package score

import (
	"fmt"
	"sync"

	"simonwaldherr.de/go/golibs/xmath"
)

var (
	criteriasMu sync.RWMutex
	criterias   = make(map[string]CriteriaEvaluator)
)

// Register makes a criteria available by the provided name.
// If Register is called twice with the same name or if criteria is nil, it panics.
func Register(name string, criteria CriteriaEvaluator) {
	criteriasMu.Lock()
	defer criteriasMu.Unlock()

	if criteria == nil {
		panic("score: Register criteria is nil")
	}
	if _, dup := criterias[name]; dup {
		panic("score: Register called twice for criteria " + name)
	}

	criterias[name] = criteria
}

// Criterias returns a list of criteria evaluators
func Criterias() map[string]CriteriaEvaluator {
	criteriasMu.RLock()
	defer criteriasMu.RUnlock()

	return criterias
}

// Messages returns a list of all criterias messages
func Messages() []string {
	criteriasMu.RLock()
	defer criteriasMu.RUnlock()

	m := []string{}

	for n, c := range criterias {
		m = append(m, fmt.Sprintf(
			"%s: %s [score = %.2f]",
			n, c.Message(),
			c.Score(),
		))
	}

	return m
}

// Weights returns a list of all criterias weights
func Weights() []float64 {
	criteriasMu.RLock()
	defer criteriasMu.RUnlock()

	w := []float64{}

	for _, c := range criterias {
		w = append(w, c.Weight())
	}

	factor := xmath.Sum(w) / 100
	for i := range w {
		w[i] /= factor
	}

	return w
}

// Values returns a list of all criterias scores
func Values() []float64 {
	criteriasMu.RLock()
	defer criteriasMu.RUnlock()

	s := []float64{}

	for _, c := range criterias {
		s = append(s, c.Score())
	}

	return s
}

// Process triggers criterias evaluation
func Process(params map[string]interface{}) float64 {
	criteriasMu.RLock()
	defer criteriasMu.RUnlock()

	for _, c := range criterias {
		c.Calculate(params)
	}

	// Loop each criterias, calculating the overall score
	s := Values()
	w := Weights()

	sw, avg := 0.0, 0.0
	for i, v := range s {
		sw += w[i]
		avg += float64(v) * w[i]
	}

	avg /= sw

	return avg
}
