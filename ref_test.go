package erx

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRet(t *testing.T) {
	var out error
	e, ret := ErrorRef(&out)

	assert.False(t, ret(nil))
	assert.Nil(t, out)

	assert.False(t, ret(e))
	assert.Nil(t, out)

	assert.False(t, ret(e, e))
	assert.Nil(t, out)

	assert.True(t, ret(fmt.Errorf("error-1")))
	assert.EqualError(t, out, "error-1")

	assert.True(t, ret(e))
	assert.EqualError(t, out, "error-1")

	assert.False(t, ret(nil))
	assert.EqualError(t, out, "error-1")

	assert.True(t, ret(e, fmt.Errorf("error-2")))
	assert.EqualError(t, out, "error-1")

	assert.True(t, ret(nil, fmt.Errorf("error-3")))
	assert.EqualError(t, out, "error-3")

	assert.True(t, ret(e, fmt.Errorf("error-4"), fmt.Errorf("error-5")))
	assert.EqualError(t, out, "error-3")

	assert.Len(t, e.suppressed, 3)
	assert.EqualError(t, e.suppressed[0], "error-2")
	assert.EqualError(t, e.suppressed[1], "error-4")
	assert.EqualError(t, e.suppressed[2], "error-5")
}

func TestError(t *testing.T) {
	var out error
	e, ret := ErrorRef(&out)

	assert.Equal(t, "<nil>", e.Error())
	ret(fmt.Errorf("error-1"))
	assert.Equal(t, "error-1", e.Error())
}

func TestWrapFmtErrorw(t *testing.T) {
	var out error
	e, ret := ErrorRef(&out)

	e.WrapFmtErrorw("test")
	assert.Nil(t, out)

	internalErr := fmt.Errorf("error-1")

	ret(internalErr)
	e.WrapFmtErrorw("test")
	assert.EqualError(t, out, "test: error-1")

	assert.True(t, errors.Is(out, internalErr))
}

func TestWrapFmtErrorf(t *testing.T) {
	var out error
	e, ret := ErrorRef(&out)

	e.WrapFmtErrorf("test: %w", OriginalErr)
	assert.Nil(t, out)

	internalErr := fmt.Errorf("error-1")

	ret(internalErr)
	e.WrapFmtErrorf("test: %w", OriginalErr)
	assert.EqualError(t, out, "test: error-1")

	assert.True(t, errors.Is(out, internalErr))

	assert.EqualError(t, OriginalErr, "<original-error-placeholder>")
}

func TestWrap(t *testing.T) {
	var out error
	e, ret := ErrorRef(&out)

	e.Wrap(func(err error) error { return fmt.Errorf("test: %w", err) })
	assert.Nil(t, out)

	internalErr := fmt.Errorf("error-1")

	ret(internalErr)
	e.Wrap(func(err error) error { return fmt.Errorf("test: %w", err) })
	assert.EqualError(t, out, "test: error-1")

	assert.True(t, errors.Is(out, internalErr))

	assert.EqualError(t, OriginalErr, "<original-error-placeholder>")
}

func TestCheckErr(t *testing.T) {
	var out error
	e, _ := ErrorRef(&out)

	e.CheckErr(func() error { return nil })
	assert.Nil(t, out)

	e.CheckErr(func() error { return fmt.Errorf("error-1") })
	assert.EqualError(t, out, "error-1")

	e.CheckErr(func() error { return fmt.Errorf("error-2") })
	assert.EqualError(t, out, "error-1")

	assert.Len(t, e.suppressed, 1)
	assert.EqualError(t, e.suppressed[0], "error-2")

	out = nil

	CheckErr(&out, func() error { return nil })
	assert.Nil(t, out)

	CheckErr(&out, func() error { return fmt.Errorf("error-1") })
	assert.EqualError(t, out, "error-1")

	CheckErr(&out, func() error { return fmt.Errorf("error-2") })
	assert.EqualError(t, out, "error-1")
}

func TestIgnoreErr(t *testing.T) {
	var out error
	e, _ := ErrorRef(&out)

	e.IgnoreErr(func() error { return fmt.Errorf("error-1") })
	assert.Nil(t, out)

	assert.Len(t, e.ignored, 1)
	assert.EqualError(t, e.ignored[0], "error-1")

	IgnoreErr(func() error { return fmt.Errorf("error-1") })
}

func TestOnError(t *testing.T) {
	var out error
	e, _ := ErrorRef(&out)

	called := false

	e.OnError(func() { called = true })
	assert.False(t, called)

	OnError(&out, func() { called = true })
	assert.False(t, called)

	out = fmt.Errorf("error-1")

	e.OnError(func() { called = true })
	assert.True(t, called)

	called = false

	OnError(&out, func() { called = true })
	assert.True(t, called)
}

func TestOnErrorOrPanic_Error(t *testing.T) {
	var out error
	e, _ := ErrorRef(&out)

	called := false

	e.OnErrorOrPanic(func() { called = true })
	assert.False(t, called)

	OnErrorOrPanic(&out, func() { called = true })
	assert.False(t, called)

	out = fmt.Errorf("error-1")

	e.OnErrorOrPanic(func() { called = true })
	assert.True(t, called)

	called = false

	OnErrorOrPanic(&out, func() { called = true })
	assert.True(t, called)
}

func TestOnPanic(t *testing.T) {
	var out error
	e, _ := ErrorRef(&out)

	called := false

	fn := func(panics bool) {
		defer e.OnPanic(func() { called = true })
		if panics {
			panic("panic")
		}
	}
	assert.NotPanics(t, func() { fn(false) })
	assert.False(t, called)

	assert.Panics(t, func() { fn(true) })
	assert.True(t, called)

	called = false
	fn2 := func(panics bool) {
		defer OnPanic(func() { called = true })
		if panics {
			panic("panic")
		}
	}
	assert.NotPanics(t, func() { fn2(false) })
	assert.False(t, called)

	assert.Panics(t, func() { fn2(true) })
	assert.True(t, called)
}

func TestOnErrorOrPanic_Panic(t *testing.T) {
	var out error
	e, _ := ErrorRef(&out)

	called := false

	fn := func(panics bool) {
		defer e.OnErrorOrPanic(func() { called = true })
		if panics {
			panic("panic")
		}
	}
	assert.NotPanics(t, func() { fn(false) })
	assert.False(t, called)

	assert.Panics(t, func() { fn(true) })
	assert.True(t, called)

	called = false
	fn2 := func(panics bool) {
		defer OnErrorOrPanic(&out, func() { called = true })
		if panics {
			panic("panic")
		}
	}
	assert.NotPanics(t, func() { fn2(false) })
	assert.False(t, called)

	assert.Panics(t, func() { fn2(true) })
	assert.True(t, called)
}

func TestOnSuccess(t *testing.T) {
	var out error
	e, _ := ErrorRef(&out)

	called := false
	e.OnSuccess(func() { called = true })
	assert.True(t, called)

	called = false
	OnSuccess(&out, func() { called = true })
	assert.True(t, called)

	out = fmt.Errorf("error-1")

	called = false
	e.OnSuccess(func() { called = true })
	assert.False(t, called)

	OnSuccess(&out, func() { called = true })
	assert.False(t, called)

	out = nil
	called = false
	fn1 := func() {
		defer e.OnSuccess(func() { called = true })
		panic("panic")
	}
	assert.Panics(t, func() { fn1() })
	assert.False(t, called)

	fn2 := func() {
		defer OnSuccess(&out, func() { called = true })
		panic("panic")
	}
	assert.Panics(t, func() { fn2() })
	assert.False(t, called)
}
