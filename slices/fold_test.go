package slices_test

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unix1/gostdx/slices"
)

type testCase[Telem any, Tacc any] struct {
	name     string
	list     []Telem
	acc      Tacc
	fn       func(Telem, Tacc) Tacc
	expected Tacc
}

type testCaseC[Telem any, Tacc any] struct {
	name     string
	list     []Telem
	acc      *Tacc
	fn       func(Telem, *Tacc) *Tacc
	expected *Tacc
}

type testCaseWithConcurrency[Telem any, Tacc any] struct {
	testCase    testCaseC[Telem, Tacc]
	concurrency int
}

func TestFold(t *testing.T) {
	testCases := []testCase[int, int]{
		{
			"simple sum",
			[]int{1, 2, 3},
			0,
			func(e int, acc int) int { return acc + e },
			6,
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := slices.Fold(tt.fn, tt.acc, tt.list)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func newFoldSumTestCase(slow bool) testCaseC[int64, int64] {
	return testCaseC[int64, int64]{
		"simple sum",
		[]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		ptr(int64(0)),
		func(e int64, acc *int64) *int64 {
			if slow {
				time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
			}
			atomic.AddInt64(acc, e)
			return acc
		},
		ptr(int64(55)),
	}
}

func TestFoldC(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testCases := []testCaseWithConcurrency[int64, int64]{
		{newFoldSumTestCase(false), 1},
		{newFoldSumTestCase(false), 3},
		{newFoldSumTestCase(false), 10},
		{newFoldSumTestCase(true), 3},
		{newFoldSumTestCase(true), 10},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(fmt.Sprintf("%s (concurrency: %d)", tt.testCase.name, tt.concurrency), func(t *testing.T) {
			actual := slices.FoldC(tt.testCase.fn, tt.testCase.acc, tt.testCase.list, tt.concurrency)
			assert.Equal(t, tt.testCase.expected, actual)
		})
	}
}

type tuple struct {
	v1 string
	v2 string
}

type mapAcc struct {
	sync.Mutex
	m map[string]string
}

func newFoldTuplesToMapTestCase() testCaseC[tuple, mapAcc] {
	return testCaseC[tuple, mapAcc]{
		"simple fold tuples to map, concurrency 1",
		[]tuple{{"k1", "v1"}, {"k2", "v2"}, {"k3", "v3"}, {"k4", "v4"}, {"k5", "v5"}, {"k6", "v6"}, {"k7", "v7"}, {"k8", "v8"}, {"k9", "v9"}, {"k10", "v10"}},
		&mapAcc{m: map[string]string{}},
		func(e tuple, acc *mapAcc) *mapAcc {
			acc.Lock()
			defer acc.Unlock()
			acc.m[e.v1] = e.v2
			return acc
		},
		&mapAcc{m: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3", "k4": "v4", "k5": "v5", "k6": "v6", "k7": "v7", "k8": "v8", "k9": "v9", "k10": "v10"}},
	}
}

func TestFoldCMap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testCases := []testCaseWithConcurrency[tuple, mapAcc]{
		{newFoldTuplesToMapTestCase(), 1},
		{newFoldTuplesToMapTestCase(), 3},
		{newFoldTuplesToMapTestCase(), 10},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(fmt.Sprintf("%s (concurrency: %d)", tt.testCase.name, tt.concurrency), func(t *testing.T) {
			actual := slices.FoldC(tt.testCase.fn, tt.testCase.acc, tt.testCase.list, tt.concurrency)
			assert.Equal(t, tt.testCase.expected, actual)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
