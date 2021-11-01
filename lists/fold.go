package lists

import (
	"sync"
)

// Fold folds (aka reduces) a list of values to a single value using the provided function and the
// starting accumulator.
// It does so by calling the `fn` function on successive elements of the list. The `fn` function
// receives the element of the list and the accumulator. It must return the new accumulator after
// processing the element. The Fold function returns the final accumulator.
func Fold[T any, Tacc any](fn func(elem T, acc Tacc) Tacc, acc0 Tacc, list []T) Tacc {
	acc := acc0
	for _, item := range list {
		acc = fn(item, acc)
	}
	return acc
}

// FoldC is a concurrent version of Fold. It takes an additional concurrency argument that spreads
// the `fn` calls over specified number of goroutines. If specified concurrency is greater than the
// length of the list, the additional goroutines will not be started.
// Keep in mind that concurrent multiple goroutines will concurrently access the accumulator.
func FoldC[T any, Tacc any](fn func(elem T, acc *Tacc) *Tacc, acc *Tacc, list []T, concurrency int) *Tacc {
	var wg sync.WaitGroup
	if concurrency > len(list) {
		concurrency = len(list)
	}
	ch := make(chan T, concurrency)
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for {
				e, ok := <-ch
				if !ok {
					return
				}
				acc = fn(e, acc)
			}
		}()
	}
	for _, item := range list {
		ch <- item
	}
	close(ch)
	wg.Wait()
	return acc
}
