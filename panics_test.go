package erx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPanicIf(t *testing.T) {
	assert.NotPanics(t, func() { PanicIf(nil) })
	assert.PanicsWithError(t, "error-1", func() {
		PanicIf(fmt.Errorf("error-1"))
	})
}

func TestPanicError(t *testing.T) {
	assert.EqualError(t, &PanicError{PanicObj: "123"}, "panic: 123")
}

func TestErrFromPanic(t *testing.T) {
	var errLogger *ErrLogger
	fn := func(panicObj interface{}, err error) (out error) {
		errLogger = ErrorLogger(&out)
		defer ErrFromPanic(&out, errLogger)

		out = err
		if panicObj != nil {
			panic(panicObj)
		}

		return err
	}

	assert.Nil(t, fn(nil, nil))
	assert.EqualError(t, fn(fmt.Errorf("error"), nil), "error")
	assert.EqualError(t, fn("error", nil), "panic: error")

	assert.EqualError(t, fn("error", fmt.Errorf("suppressed")), "panic: error")
	assert.Len(t, errLogger.suppressed, 1)
	assert.EqualError(t, errLogger.suppressed[0], "suppressed")
}
