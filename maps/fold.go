package maps

import (
	"sync"
)

type msg[Tkey comparable, Tval any] struct {
	key Tkey
	val Tval
}

// Fold folds (aka reduces) a map of values to a single value using the provided function and the
// starting accumulator.
// It does so by calling the `fn` function on successive elements of the map. The `fn` function
// receives each key-value pair from the map and the accumulator. It must return the new accumulator
// after processing the key-value pair. The Fold function returns the final accumulator.
func Fold[Tkey comparable, Tval any, Tacc any](fn func(key Tkey, val Tval, acc Tacc) Tacc, acc0 Tacc, m map[Tkey]Tval) Tacc {
	acc := acc0
	for key, val := range m {
		acc = fn(key, val, acc)
	}
	return acc
}

// FoldC is a concurrent version of Fold. It takes an additional concurrency argument that spreads
// the `fn` calls over specified number of goroutines. If specified concurrency is greater than the
// length of the map, the additional goroutines will not be started.
// Keep in mind that multiple goroutines will concurrently access the accumulator; and the
// accumulator is a pointer.
func FoldC[Tkey comparable, Tval any, Tacc any](fn func(key Tkey, val Tval, acc *Tacc) *Tacc, acc *Tacc, m map[Tkey]Tval, concurrency int) *Tacc {
	var wg sync.WaitGroup
	if concurrency > len(m) {
		concurrency = len(m)
	}
	ch := make(chan msg[Tkey, Tval], concurrency)
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for {
				e, ok := <-ch
				if !ok {
					return
				}
				acc = fn(e.key, e.val, acc)
			}
		}()
	}
	for key, val := range m {
		ch <- msg[Tkey, Tval]{key, val}
	}
	close(ch)
	wg.Wait()
	return acc
}
