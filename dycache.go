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
	if b.Exist(v) {
		b.Purge(v)
		// l := b.cache.Len()
		// b.best = (b.best*80 + b.size*20) / 100
		// fmt.Printf("miss, cache-size:%d, miss-length:%d\n", l, b.size)
		b.Summary(b.size)
		return
	}
	b.Set(v)
}

type String2Int func(string) int

func s2i(s string) int {
	r := 0
	for _, it := range s {
		r = r*10 + int(it)
	}
	return r
}
