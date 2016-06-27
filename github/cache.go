package github

import "strings"

func getCached(key ...string) (interface{}, bool) {
	k := strings.Join(key, "-")
	return cache.Get(k)
}

func saveCache(value interface{}, key ...string) {
	k := strings.Join(key, "-")
	cache.Add(k, value)
}
