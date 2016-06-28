package github

import (
	"reflect"
	"testing"

	"github.com/hashicorp/golang-lru"
)

func TestCacheInitialised(t *testing.T) {
	var err error
	cache, err = lru.NewARC(CacheSize)
	if err != nil {
		t.Errorf("The cache could not be initialised")
	}
}

func TestCacheSaved(t *testing.T) {
	saveCache(map[string]int{"foo": 1}, "foo", "bar")
	if cache.Len() != 1 {
		t.Error("The item has not been cached")
	}
}

func TestCacheLoaded(t *testing.T) {
	val, exists := getCached("foo", "bar")
	if !exists {
		t.Error("The item should be in cache")
	}
	if !reflect.DeepEqual(val, map[string]int{"foo": 1}) {
		t.Error("The item in cache is not the same")
	}
}
