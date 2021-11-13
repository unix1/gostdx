package maps_test

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unix1/gostdx/maps"
)

type testCase[Tkey comparable, Tval any, Tacc any] struct {
	name     string
	m        map[Tkey]Tval
	acc      Tacc
	fn       func(Tkey, Tval, Tacc) Tacc
	expected Tacc
}

type testCaseC[Tkey comparable, Tval any, Tacc any] struct {
	name     string
	m        map[Tkey]Tval
	acc      *Tacc
	fn       func(Tkey, Tval, *Tacc) *Tacc
	expected *Tacc
}

type testCaseWithConcurrency[Tkey comparable, Tval any, Tacc any] struct {
	testCase    testCaseC[Tkey, Tval, Tacc]
	concurrency int
}

func TestFold(t *testing.T) {
	testCases := []testCase[int, int, int]{
		{
			"simple sum of k*v",
			map[int]int{1: 10, 2: 20, 3: 30},
			0,
			func(k int, v int, acc int) int { return acc + k*v },
			140,
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := maps.Fold(tt.fn, tt.acc, tt.m)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func newFoldSumTestCase(slow bool) testCaseC[int64, int64, int64] {
	return testCaseC[int64, int64, int64]{
		"simple sum",
		map[int64]int64{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10},
		ptr(int64(0)),
		func(k int64, v int64, acc *int64) *int64 {
			if slow {
				time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
			}
			atomic.AddInt64(acc, k*v)
			return acc
		},
		ptr(int64(385)),
	}
}

func TestFoldC(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testCases := []testCaseWithConcurrency[int64, int64, int64]{
		{newFoldSumTestCase(false), 1},
		{newFoldSumTestCase(false), 3},
		{newFoldSumTestCase(false), 10},
		{newFoldSumTestCase(true), 3},
		{newFoldSumTestCase(true), 10},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(fmt.Sprintf("%s (concurrency: %d)", tt.testCase.name, tt.concurrency), func(t *testing.T) {
			actual := maps.FoldC(tt.testCase.fn, tt.testCase.acc, tt.testCase.m, tt.concurrency)
			assert.Equal(t, tt.testCase.expected, actual)
		})
	}
}

type tuple struct {
	v1 string
	v2 string
}

type tuples struct {
	sync.Mutex
	t []tuple
}

func newFoldTuplesToMapTestCase() testCaseC[string, string, tuples] {
	return testCaseC[string, string, tuples]{
		"simple fold tuples to map, concurrency 1",
		map[string]string{"k1": "v1", "k2": "v2", "k3": "v3", "k4": "v4", "k5": "v5", "k6": "v6", "k7": "v7", "k8": "v8", "k9": "v9", "k10": "v10"},
		&tuples{t: []tuple{}},
		func(k, v string, acc *tuples) *tuples {
			acc.Lock()
			defer acc.Unlock()
			acc.t = append(acc.t, tuple{v1: k, v2: v})
			return acc
		},
		&tuples{t: []tuple{{"k1", "v1"}, {"k2", "v2"}, {"k3", "v3"}, {"k4", "v4"}, {"k5", "v5"}, {"k6", "v6"}, {"k7", "v7"}, {"k8", "v8"}, {"k9", "v9"}, {"k10", "v10"}}},
	}
}

func TestFoldCMap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testCases := []testCaseWithConcurrency[string, string, tuples]{
		{newFoldTuplesToMapTestCase(), 1},
		{newFoldTuplesToMapTestCase(), 3},
		{newFoldTuplesToMapTestCase(), 10},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(fmt.Sprintf("%s (concurrency: %d)", tt.testCase.name, tt.concurrency), func(t *testing.T) {
			actual := maps.FoldC(tt.testCase.fn, tt.testCase.acc, tt.testCase.m, tt.concurrency)
			assert.ElementsMatch(t, tt.testCase.expected.t, actual.t)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
