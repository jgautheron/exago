package github

import "strings"

func (g GitHub) getCached(key ...string) (interface{}, bool) {
	g.RLock()
	defer g.RUnlock()
	k := strings.Join(key, "-")
	return g.cache.Get(k)
}

func (g GitHub) saveCache(value interface{}, key ...string) {
	g.Lock()
	defer g.Unlock()
	k := strings.Join(key, "-")
	g.cache.Add(k, value)
}
