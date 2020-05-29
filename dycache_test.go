package gcache

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/golang/groupcache/lru"
)

var lock sync.Mutex

func setIfNotEx(cache *lru.Cache, key lru.Key, value interface{}) bool {
	lock.Lock()
	defer lock.Unlock()
	_, ex := cache.Get(key)
	if ex {
		return true
	}
	cache.Add(key, value)
	return false
}

func TestOnEvicted(t *testing.T) {
	cache := lru.New(799)
	b := WarpCache(cache, nil)
	_ = b

	r := rand.New(rand.NewSource(time.Now().Unix()))

	hit := 0
	for i := 0; i < 600000; i++ {
		a := r.Intn(1000)
		if setIfNotEx(cache, a, i) {
			hit++
		}
	}
	t.Logf("hit:%d", hit)
}
