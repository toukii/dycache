package gcache

import (
	"fmt"
	"github.com/golang/groupcache/lru"
)

type Bitmap struct {
	bits  [3][2048]bool
	cache *lru.Cache
	size  int
	best  int

	set     chan int
	summary chan struct{}
	vs      []int

	String2Int
}

var primes [3]uint

func init() {
	primes = [3]uint{1, 1, 1}
}

func (b *Bitmap) Exist(v int) bool {
	for i, prime := range primes {
		v = int(v << prime % 2048)
		if !b.bits[i][v] {
			return false
		}
	}
	return true
}

func (b *Bitmap) Set(v int) {
	b.size++
	for i, prime := range primes {
		v = int(v << prime % 2048)
		b.bits[i][v] = true
	}
}

func (b *Bitmap) Purge(v int) {
	b.size--
	for i, prime := range primes {
		v = int(v << prime % 2048)
		b.bits[i][v] = false
	}
}

func (b *Bitmap) Summary(l int) {
	if b.set == nil {
		b.set = make(chan int, 10)
	}
	if b.vs == nil {
		b.vs = make([]int, 10)
	}
	if b.summary == nil {
		b.summary = make(chan struct{}, 2)
	}

	select {
	case b.set <- l:
		if len(b.set) >= 10 {
			// b.summary <- struct{}{}
			size := len(b.set)
			for i := 0; i < size; i++ {
				b.vs[i] = <-b.set
			}
			fmt.Printf("avg:%d\n", avg(b.vs[:size]))
		}
		// case <-b.summary:

	}

}

func avg(vs []int) int {
	sum := 0
	for _, v := range vs {
		sum += v
	}
	return sum * 8 / len(vs) / 10
}
