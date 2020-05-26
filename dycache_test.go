package gcache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/groupcache/lru"
)

func setIfNotEx(cache *lru.Cache, key lru.Key, value interface{}) bool {
	// fmt.Printf("%+v ", key)
	_, ex := cache.Get(key)
	if ex {
		return true
	}
	cache.Add(key, value)
	return false
}

func TestOnEvicted(t *testing.T) {
	cache := lru.New(80)
	b := WarpCache(cache, nil)
	_ = b

	for i := 0; i < 150; i++ {
		cache.Add(i, i)
		// fmt.Printf("%d ", i)
	}
	fmt.Println("=====")

	r := rand.New(rand.NewSource(time.Now().Unix()))

	hit := 0
	for i := 0; i < 1500; i++ {
		a := r.Intn(100)
		if setIfNotEx(cache, a, i) {
			hit++
		}
	}
	t.Logf("hit:%d", hit)
	hit = 0
	for i := 0; i < 1500; i++ {
		a := r.Intn(100)
		if setIfNotEx(cache, a, i) {
			hit++
		}
	}
	t.Logf("hit:%d", hit)
}
