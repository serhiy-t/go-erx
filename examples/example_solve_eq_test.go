package examples

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSquareEq(t *testing.T) {
	x1, x2, err := squareEq("/solve-quadratic/a/2/b/4/c/6")
	assert.EqualError(t, err, "expected: D >= 0; actual D = -32.000000")
	assert.Equal(t, 0.0, x1)
	assert.Equal(t, 0.0, x2)

	x1, x2, err = squareEq("/solve-quadratic/a/2.0/b/5.0/c/-3.0")
	assert.Nil(t, err)
	assert.Equal(t, -3.0, x1)
	assert.Equal(t, 0.5, x2)

	x1, x2, err = squareEq("/solve-quadratic/a/2.0/z/5.0/c/-3.0")
	assert.EqualError(t, err, "expected 'b' next, but got: z/5.0/c/-3.0")
	assert.Equal(t, 0.0, x1)
	assert.Equal(t, 0.0, x2)
}
