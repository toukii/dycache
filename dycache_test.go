package gcache

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/golang/groupcache/lru"
)

var lock sync.Mutex
var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

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
	cache := lru.New(90)
	WarpCache(cache, nil, time.Millisecond*20000)

	// sample := 6000
	sample := 6000000

	data := simulateColdData()
	// data := simulateRandData()
	hit := 0
	m := make(map[int]bool, 100)
	for i := 0; i < sample; i++ {
		a := r.Intn(100)
		// if a > 80 {
		a += <-data
		// }
		m[a] = true
		if setIfNotEx(cache, a, i) {
			hit++
		}
	}

	t.Logf("hit:%d, rate:%.2f%%", hit, float64(hit)*100/float64(sample))
	t.Logf("data-range:%d\n%+v", len(m), m)
}

// cold-data: 100ms
func simulateColdData() chan int {
	// ticker := time.NewTicker()
	timer := time.NewTimer(time.Millisecond * 200)
	now := time.Now().UnixNano()
	ch := make(chan int, 2)
	go func(now *int64) {
		for {
			<-timer.C
			*now += 50
		}
	}(&now)
	go func() {
		for {
			// fmt.Println(now)
			if now%500 > 100 {
				ch <- 100
				fmt.Print(".")
			} else {
				ch <- 0
			}
		}
	}()
	return ch
}

// rand-data: plus 100
func simulateRandData() chan int {
	ch := make(chan int, 2)
	go func() {
		for {
			i := time.Now().UnixNano()
			if i%2 == 1 {
				ch <- 100
			} else {
				ch <- 0
			}
		}
	}()
	return ch
}
