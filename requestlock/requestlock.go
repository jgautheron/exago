package requestlock

import "strings"

var (
	lock = map[string]bool{}
)

func Has(args ...string) bool {
	p := strings.Join(args, "-")
	_, found := lock[p]
	return found
}

func Add(args ...string) {
	p := strings.Join(args, "-")
	lock[p] = true
}

func Remove(args ...string) {
	p := strings.Join(args, "-")
	delete(lock, p)
}
