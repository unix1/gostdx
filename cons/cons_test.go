package cons_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/unix1/gostdx/cons"
)

func TestNewEmpty(t *testing.T) {
	c := Cons([]int{}...)
	assert.Nil(t, c)
}

func TestNewCons(t *testing.T) {
	c := Cons(1, 2, 3)
	assert.Equal(t, 1, Car(c))
	assert.Equal(t, 2, Car(Cdr(c)))
	assert.Equal(t, 3, Car(Cdr(Cdr(c))))
	assert.Nil(t, Cdr(Cdr(Cdr(c))))
}
