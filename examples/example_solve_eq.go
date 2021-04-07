package examples

import (
	"github.com/serhiy-t/go-erx"
	"math"
)

// url: /solve-quadratic/a/2/b/4/c/6
func squareEq(url string) (x1, x2 float64, out error) {
	defer erx.ErrFromPanic(&out)

	consumeKeyOrPanic(&url, "solve-quadratic")
	consumeKeyOrPanic(&url, "a")
	a := atofOrPanic(consumeValueOrPanic(&url))
	consumeKeyOrPanic(&url, "b")
	b := atofOrPanic(consumeValueOrPanic(&url))
	consumeKeyOrPanic(&url, "c")
	c := atofOrPanic(consumeValueOrPanic(&url))

	d := b*b - 4*a*c
	erx.PanicIf(erx.AssertErr(d >= 0.0, "expected: D >= 0; actual D = %f", d))

	dSqrt := math.Sqrt(b*b - 4*a*c)

	x1 = (-b - dSqrt) / (2 * a)
	x2 = (-b + dSqrt) / (2 * a)

	return
}
