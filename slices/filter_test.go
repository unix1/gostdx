package slices_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unix1/gostdx/slices"
)

type filterTestCase[T any] struct {
	name     string
	list     []T
	fn       func(elem T) bool
	expected []T
}

type filterTestCaseWithConcurrency[T any] struct {
	testCase    filterTestCase[T]
	concurrency int
}

func TestFilter(t *testing.T) {
	testCases := []filterTestCase[int]{
		{
			"empty list",
			[]int{},
			func(e int) bool { return false },
			[]int{},
		},
		{
			"none match",
			[]int{1, 2, 3, 4, 5},
			func(e int) bool { return false },
			[]int{},
		},
		{
			"even only",
			[]int{1, 2, 3, 4, 5},
			func(e int) bool { return e%2 == 0 },
			[]int{2, 4},
		},
		{
			"all match",
			[]int{1, 2, 3, 4, 5},
			func(e int) bool { return true },
			[]int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Filter(tt.fn, tt.list)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func newFilterEvenTestCase(slow bool) filterTestCase[int] {
	return filterTestCase[int]{
		"even only",
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		func(e int) bool {
			if slow {
				time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
			}
			return e%2 == 0
		},
		[]int{2, 4, 6, 8, 10},
	}
}

func TestFilterC(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testCases := []filterTestCaseWithConcurrency[int]{
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
			actual := slices.FilterC(tt.testCase.fn, tt.testCase.list, tt.concurrency)
			assert.ElementsMatch(t, tt.testCase.expected, actual)
		})
	}
}
