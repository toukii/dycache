package gcache

import (
	"fmt"
	"time"

	"github.com/golang/groupcache/lru"
	// "github.com/willf/bitset"
)

const (
	mask = 10240
)

type Bitmap struct {
	bits  [3][mask]int64
	cache *lru.Cache
	size  int
	best  int
	ex    int64

	String2Int
}

var primes [3]uint

func init() {
	primes = [3]uint{1, 1, 1}
}

func NewBitmap(cache *lru.Cache, ex time.Duration) *Bitmap {
	return &Bitmap{
		cache: cache,
		ex:    ex.Nanoseconds(),
	}
}

func (b *Bitmap) bitmap(v int) [3]int {
	vs := [3]int{}
	for i, prime := range primes {
		v = int(v << prime % mask)
		vs[i] = v
	}
	return vs
}

func setodd(v int64) int64 {
	if v%2 == 1 {
		return v
	}
	return v - 1
}

func (b *Bitmap) crease(vs [3]int) bool {
	now := time.Now().UnixNano()
	exp := b.bits[0][vs[0]] > 0 && now-b.bits[0][vs[0]] > b.ex // 是否过期
	equal := b.bits[0][vs[0]] == b.bits[1][vs[1]] && b.bits[1][vs[1]] == b.bits[2][vs[2]]
	odd := setodd(now)
	if exp {
		fmt.Println("exp ", b.bits[0][vs[0]], b.bits[1][vs[1]], b.bits[2][vs[2]])
	}
	if !exp && equal {
		isodd := b.bits[0][vs[0]]%2 == 1 // the one before is odd
		if isodd {
			b.size--
			b.bits[0][vs[0]], b.bits[1][vs[1]], b.bits[2][vs[2]] = odd+1, odd+1, odd+1
		} else {
			b.size++
			b.bits[0][vs[0]], b.bits[1][vs[1]], b.bits[2][vs[2]] = odd, odd, odd
		}
	} else if equal && exp {
		b.bits[0][vs[0]], b.bits[1][vs[1]], b.bits[2][vs[2]] = odd+1, odd+1, odd+1 // set old+1, next will plus b.size
		// b.bits[0][vs[0]], b.bits[1][vs[1]], b.bits[2][vs[2]] = 0, 0, 0 // set old+1, next will plus b.size
		// b.size--
	} else {
		fmt.Println(b.bits[0][vs[0]], b.bits[1][vs[1]], b.bits[2][vs[2]], vs[0], vs[1], vs[2])
		return false
	}
	return true
}

func (b *Bitmap) Set(v int) {
	vs := b.bitmap(v)
	if !b.crease(vs) {
		fmt.Println(v)
		b.bits[0][vs[0]] = 0
		b.bits[1][vs[1]] = 0
		b.bits[2][vs[2]] = 0
		b.size++
	}
}

func (b *Bitmap) ExistPurge(v int) bool {
	vs := b.bitmap(v)

	if !b.crease(vs) {
		return false
	}

	b.best++

	b.Summary(b.size)

	return true
}

func (b *Bitmap) Exist(v int) bool {
	for i, prime := range primes {
		v = int(v << prime % mask)
		if b.bits[i][v] == 0 {
			return false
		}
	}
	return true
}

func (b *Bitmap) Purge(v int) {
	b.size--
	for i, prime := range primes {
		v = int(v << prime % mask)
		b.bits[i][v] = 0
	}
}

func (b *Bitmap) Summary(l int) {
	// b.best = (b.best*2 + l*8) * 8 / 100
	// b.best = l
	// le, be := b.cache.Len(), b.size>>1
	le, be := b.cache.Len(), b.size<<1
	rate := float32(le) / float32(be)
	// if be > le {
	if be > le && be < le*11/10 {
		// b.cache.MaxEntries = be
	}
	// return
	fmt.Printf("\t%d ==> %d (%.2f)\t", le, be, rate)
}

func avg(vs []int) int {
	sum := 0
	for _, v := range vs {
		sum += v
	}
	return sum * 8 / len(vs) / 10
}
