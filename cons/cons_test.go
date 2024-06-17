package cons_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/unix1/gostdx/cons"
)

func TestNewEmpty(t *testing.T) {
	c := List([]int{}...)
	assert.Nil(t, c)
}

func TestNewList(t *testing.T) {
	c := List(1, 2, 3)
	assert.Equal(t, 1, Car(c))
	assert.Equal(t, 2, Car(Cdr(c)))
	assert.Equal(t, 3, Car(Cdr(Cdr(c))))
	assert.Nil(t, Cdr(Cdr(Cdr(c))))
}

func TestIterEmpty(t *testing.T) {
	c := List[int]()
	var want []int
	var got []int
	for v := range Each(c) {
		got = append(got, v)
	}
	assert.EqualValues(t, want, got)
}

func TestIterEach(t *testing.T) {
	c := List(1, 2, 3)
	want := []int{1, 2, 3}
	var got []int
	for v := range Each(c) {
		got = append(got, v)
	}
	fmt.Println("got", got)
	assert.EqualValues(t, want, got)
}
