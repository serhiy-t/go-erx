package examples

import (
	"github.com/serhiy-t/go-erx"
	"math"
)

// url: /solve-quadratic/a/2/b/4/c/6
func squareEq(url string) (x1, x2 float64, out error) {
	defer erx.ErrFromPanic(&out)

	panicsConsumeKey(&url, "solve-quadratic")
	panicsConsumeKey(&url, "a")
	a := panicsAtof(panicsConsumeValue(&url))
	panicsConsumeKey(&url, "b")
	b := panicsAtof(panicsConsumeValue(&url))
	panicsConsumeKey(&url, "c")
	c := panicsAtof(panicsConsumeValue(&url))

	d := b*b - 4*a*c
	erx.PanicIf(erx.AssertErr(d >= 0.0, "expected: D >= 0; actual D = %f", d))

	dSqrt := math.Sqrt(b*b - 4*a*c)

	x1 = (-b - dSqrt) / (2 * a)
	x2 = (-b + dSqrt) / (2 * a)

	return
}
