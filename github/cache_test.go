package github

import (
	"reflect"
	"testing"
)

var githubTest GitHub

func init() {
	githubTest = GetInstance()
}

func TestCacheSaved(t *testing.T) {
	githubTest.saveCache(map[string]int{"foo": 1}, "foo", "bar")
	if githubTest.cache.Len() != 1 {
		t.Error("The item has not been cached")
	}
}

func TestCacheLoaded(t *testing.T) {
	val, exists := githubTest.getCached("foo", "bar")
	if !exists {
		t.Error("The item should be in cache")
	}
	if !reflect.DeepEqual(val, map[string]int{"foo": 1}) {
		t.Error("The item in cache is not the same")
	}
}
