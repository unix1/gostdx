package maps

import (
	"sync"
)

type kv[Tkey comparable, Tval any] struct {
	k Tkey
	v Tval
}

// Filter filters out the elements of the given list that don't satisfy the predicate function fn.
func Filter[Tkey comparable, Tval any](fn func(key Tkey, val Tval) bool, m map[Tkey]Tval) map[Tkey]Tval {
	result := make(map[Tkey]Tval)
	for key, val := range m {
		if fn(key, val) {
			result[key] = val
		}
	}
	return result
}

// FilterC is a concurrent version of Filter.
// The additional concurrency argument ensures at most that many concurrent goroutines are doing
// work.
func FilterC[Tkey comparable, Tval any](fn func(key Tkey, val Tval) bool, m map[Tkey]Tval, concurrency int) map[Tkey]Tval {
	var wg, resWg sync.WaitGroup
	if concurrency > len(m) {
		concurrency = len(m)
	}
	result := make(map[Tkey]Tval)
	ch := make(chan kv[Tkey, Tval], concurrency)
	resCh := make(chan kv[Tkey, Tval], concurrency)
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for {
				e, ok := <-ch
				if !ok {
					return
				}
				if fn(e.k, e.v) {
					resCh <- e
				}
			}
		}()
	}
	resWg.Add(1)
	go func() {
		defer resWg.Done()
		for {
			e, ok := <-resCh
			if !ok {
				break
			}
			result[e.k] = e.v
		}
	}()
	for key, val := range m {
		ch <- kv[Tkey, Tval]{k: key, v: val}
	}
	close(ch)
	wg.Wait()
	close(resCh)
	resWg.Wait()
	return result
}
