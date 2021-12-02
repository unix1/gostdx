package slices

import (
	"sync"
)

// Filter filters out the elements of the given list that don't satisfy the predicate function fn.
func Filter[T any](fn func(elem T) bool, list []T) []T {
	// optimistic capacity
	result := make([]T, 0, len(list))
	for _, item := range list {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// FilterC is a concurrent version of Filter.
// The additional concurrency argument ensures at most that many concurrent goroutines are doing
// work.
func FilterC[T any](fn func(elem T) bool, list []T, concurrency int) []T {
	var wg, resWg sync.WaitGroup
	if concurrency > len(list) {
		concurrency = len(list)
	}
	// optimistic capacity
	result := make([]T, 0, len(list))
	ch := make(chan T, concurrency)
	resCh := make(chan T, concurrency)
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for {
				e, ok := <-ch
				if !ok {
					return
				}
				if fn(e) {
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
			result = append(result, e)
		}
	}()
	for _, item := range list {
		ch <- item
	}
	close(ch)
	wg.Wait()
	close(resCh)
	resWg.Wait()
	return result
}
