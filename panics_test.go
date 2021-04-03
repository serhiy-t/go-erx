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
	fn := func(panicObj interface{}) (out error) {
		defer ErrFromPanic(&out)

		if panicObj != nil {
			panic(panicObj)
		}

		return nil
	}

	assert.Nil(t, fn(nil))
	assert.EqualError(t, fn(fmt.Errorf("error")), "error")
	assert.EqualError(t, fn("error"), "panic: error")
}
