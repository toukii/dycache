package gcache

import (
	"fmt"
	"github.com/golang/groupcache/lru"
)

const (
	mask = 20140
)

type Bitmap struct {
	bits  [3][mask]bool
	cache *lru.Cache
	size  int
	best  int

	String2Int
}

var primes [3]uint

func init() {
	primes = [3]uint{1, 1, 1}
}

func (b *Bitmap) bitmap(v int) [3]int {
	vs := [3]int{}
	for i, prime := range primes {
		v = int(v << prime % mask)
		vs[i] = v
	}
	return vs
}

func (b *Bitmap) crease(vs [3]int) bool {
	if b.bits[0][vs[0]] && b.bits[1][vs[1]] && b.bits[2][vs[2]] {
		b.size--
		b.bits[0][vs[0]], b.bits[1][vs[1]], b.bits[2][vs[2]] = false, false, false
	} else if !b.bits[0][vs[0]] && !b.bits[1][vs[1]] && !b.bits[2][vs[2]] {
		b.size++
		b.bits[0][vs[0]], b.bits[1][vs[1]], b.bits[2][vs[2]] = true, true, true
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
		b.bits[0][vs[0]] = true
		b.bits[1][vs[1]] = true
		b.bits[2][vs[2]] = true
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
		if b.bits[i][v] {
			return false
		}
	}
	return true
}

func (b *Bitmap) Purge(v int) {
	b.size--
	for i, prime := range primes {
		v = int(v << prime % mask)
		b.bits[i][v] = false
	}
}

func (b *Bitmap) Summary(l int) {
	// b.best = (b.best*2 + l*8) * 8 / 100
	// b.best = l
	le, be := b.cache.Len(), b.size<<1
	rate := float32(le) / float32(be)
	if be > le {
		// if be > le && be < le*11/10 {
		// b.cache.MaxEntries = be
	}
	fmt.Printf("%d ==> %d (%.2f) %d\n", le, be, rate, b.best)
}

func avg(vs []int) int {
	sum := 0
	for _, v := range vs {
		sum += v
	}
	return sum * 8 / len(vs) / 10
}
