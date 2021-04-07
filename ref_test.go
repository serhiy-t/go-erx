package erx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrapFmtErrorw(t *testing.T) {
	var out error

	WrapFmtErrorw(&out, "test")
	assert.Nil(t, out)

	out = fmt.Errorf("error-1")
	WrapFmtErrorw(&out, "test")
	assert.EqualError(t, out, "test: error-1")
}

func TestWrapFmtErrorf(t *testing.T) {
	var out error

	WrapFmtErrorf(&out, "test: %w", OriginalErr)
	assert.Nil(t, out)

	out = fmt.Errorf("error-1")

	WrapFmtErrorf(&out, "test: %w", OriginalErr)
	assert.EqualError(t, out, "test: error-1")

	assert.EqualError(t, OriginalErr, "<original-error-placeholder>")
}

func TestWrap(t *testing.T) {
	var out error

	Wrap(&out, func(err error) error { return fmt.Errorf("test: %w", err) })
	assert.Nil(t, out)

	out = fmt.Errorf("error-1")

	Wrap(&out, func(err error) error { return fmt.Errorf("test: %w", err) })
	assert.EqualError(t, out, "test: error-1")

	assert.EqualError(t, OriginalErr, "<original-error-placeholder>")
}

func TestCheckErr(t *testing.T) {
	var out error
	errLogger := ErrorLogger(&out)

	CheckErr(&out, func() error { return nil }, errLogger)
	assert.Nil(t, out)

	CheckErr(&out, func() error { return fmt.Errorf("error-1") }, errLogger)
	assert.EqualError(t, out, "error-1")

	CheckErr(&out, func() error { return fmt.Errorf("error-2") }, errLogger)
	assert.EqualError(t, out, "error-1")

	assert.Len(t, errLogger.suppressed, 1)
	assert.EqualError(t, errLogger.suppressed[0], "error-2")
}

func TestIgnoreErr(t *testing.T) {
	var out error
	errLogger := ErrorLogger(&out)

	IgnoreErr(func() error { return fmt.Errorf("error-1") }, errLogger)
	assert.Nil(t, out)

	assert.Len(t, errLogger.ignored, 1)
	assert.EqualError(t, errLogger.ignored[0], "error-1")
}

func TestOnError(t *testing.T) {
	var out error
	called := false

	OnError(&out, func() { called = true })
	assert.False(t, called)

	out = fmt.Errorf("error-1")

	called = false

	OnError(&out, func() { called = true })
	assert.True(t, called)
}

func TestOnErrorOrPanic_Error(t *testing.T) {
	var out error

	called := false

	OnErrorOrPanic(&out, func() { called = true })
	assert.False(t, called)

	out = fmt.Errorf("error-1")

	called = false

	OnErrorOrPanic(&out, func() { called = true })
	assert.True(t, called)
}

func TestOnPanic(t *testing.T) {
	called := false

	fn := func(panics bool) {
		defer OnPanic(func() { called = true })
		if panics {
			panic("panic")
		}
	}
	assert.NotPanics(t, func() { fn(false) })
	assert.False(t, called)

	assert.Panics(t, func() { fn(true) })
	assert.True(t, called)
}

func TestOnErrorOrPanic_Panic(t *testing.T) {
	var out error

	called := false
	fn := func(panics bool) {
		defer OnErrorOrPanic(&out, func() { called = true })
		if panics {
			panic("panic")
		}
	}
	assert.NotPanics(t, func() { fn(false) })
	assert.False(t, called)

	assert.Panics(t, func() { fn(true) })
	assert.True(t, called)
}

func TestOnSuccess(t *testing.T) {
	var out error

	called := false
	OnSuccess(&out, func() { called = true })
	assert.True(t, called)

	out = fmt.Errorf("error-1")

	called = false
	OnSuccess(&out, func() { called = true })
	assert.False(t, called)

	out = nil
	called = false

	fn := func() {
		defer OnSuccess(&out, func() { called = true })
		panic("panic")
	}
	assert.Panics(t, func() { fn() })
	assert.False(t, called)
}
