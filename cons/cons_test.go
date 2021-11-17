package cons_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unix1/gostdx/cons"
)

func TestNewEmpty(t *testing.T) {
	c := cons.New([]int{}...)
	assert.Nil(t, c)
}

func TestNewCons(t *testing.T) {
	c := cons.New(1, 2, 3)
	assert.Equal(t, 1, cons.Car(c))
	assert.Equal(t, 2, cons.Car(cons.Cdr(c)))
	assert.Equal(t, 3, cons.Car(cons.Cdr(cons.Cdr(c))))
	assert.Nil(t, cons.Cdr(cons.Cdr(cons.Cdr(c))))
}
