package github

import "strings"

func (g GitHub) getCached(key ...string) (interface{}, bool) {
	k := strings.Join(key, "-")
	return g.cache.Get(k)
}

func (g GitHub) saveCache(value interface{}, key ...string) {
	k := strings.Join(key, "-")
	g.cache.Add(k, value)
}
