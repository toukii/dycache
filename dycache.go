package gcache

import (
	"fmt"
	"github.com/golang/groupcache/lru"
)

func WarpCache(cache *lru.Cache, f String2Int) *Bitmap {
	b := new(Bitmap)
	b.cache = cache
	b.String2Int = f
	if f == nil {
		b.String2Int = s2i
	}
	cache.OnEvicted = b.OnEvicted
	return b
}

func (b *Bitmap) OnEvicted(key lru.Key, value interface{}) {
	v := b.String2Int(fmt.Sprintf("%+v", key))
	b.ExistPurge(v)
}

type String2Int func(string) int

func s2i(s string) int {
	r := 0
	for _, it := range s {
		r = r*10 + int(it)
	}
	return r
}
