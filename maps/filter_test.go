package maps_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unix1/gostdx/maps"
)

type filterTestCase[Tkey comparable, Tval any] struct {
	name     string
	m        map[Tkey]Tval
	fn       func(Tkey, Tval) bool
	expected map[Tkey]Tval
}

type filterTestCaseWithConcurrency[Tkey comparable, Tval any] struct {
	testCase    filterTestCase[Tkey, Tval]
	concurrency int
}

func TestFilter(t *testing.T) {
	testCases := []filterTestCase[int, int]{
		{
			"empty map",
			map[int]int{},
			func(k, v int) bool { return false },
			map[int]int{},
		},
		{
			"none match",
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5},
			func(k, v int) bool { return false },
			map[int]int{},
		},
		{
			"even only",
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5},
			func(k, v int) bool { return v%2 == 0 },
			map[int]int{2: 2, 4: 4},
		},
		{
			"all match",
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5},
			func(k, v int) bool { return true },
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := maps.Filter(tt.fn, tt.m)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func newFilterEvenTestCase(slow bool) filterTestCase[int, int] {
	return filterTestCase[int, int]{
		"even only",
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10},
		func(k, v int) bool {
			if slow {
				time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
			}
			return v%2 == 0
		},
		map[int]int{2: 2, 4: 4, 6: 6, 8: 8, 10: 10},
	}
}

func TestFilterC(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testCases := []filterTestCaseWithConcurrency[int, int]{
		{testCase: newFilterEvenTestCase(false), concurrency: 1},
		{testCase: newFilterEvenTestCase(false), concurrency: 5},
		{testCase: newFilterEvenTestCase(false), concurrency: 10},
		{testCase: newFilterEvenTestCase(true), concurrency: 1},
		{testCase: newFilterEvenTestCase(true), concurrency: 5},
		{testCase: newFilterEvenTestCase(true), concurrency: 10},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.testCase.name, func(t *testing.T) {
			actual := maps.FilterC(tt.testCase.fn, tt.testCase.m, tt.concurrency)
			assert.Equal(t, tt.testCase.expected, actual)
		})
	}
}
