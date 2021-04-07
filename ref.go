package erx

import "fmt"

func Wrap(errPtr *error, fn func(error) error) {
	err := *errPtr
	if err == nil {
		return
	}
	*errPtr = fn(err)
}

type originalErrorType struct{}

func (originalErrorType) Error() string {
	return "<original-error-placeholder>"
}

var OriginalErr error = originalErrorType{}

func WrapFmtErrorf(errPtr *error, format string, a ...interface{}) {
	if *errPtr == nil {
		return
	}
	Wrap(errPtr, func(err error) error {
		for i := range a {
			if a[i] == OriginalErr {
				a[i] = err
			}
		}
		return fmt.Errorf(format, a...)
	})
}

func WrapFmtErrorw(errPtr *error, format string, a ...interface{}) {
	if *errPtr == nil {
		return
	}
	WrapFmtErrorf(errPtr, format+": %w", append(a, OriginalErr)...)
}

func CheckErr(errPtr *error, fn func() error, options ...Option) {
	fnErr := fn()
	if fnErr != nil {
		if *errPtr == nil {
			*errPtr = fnErr
		} else {
			for _, o := range options {
				o.reportSuppressed(fnErr)
			}
		}
	}
}

func IgnoreErr(fn func() error, options ...Option) {
	ignoredErr := fn()
	if ignoredErr != nil {
		for _, o := range options {
			o.reportIgnored(ignoredErr)
		}
	}
}

func OnError(errPtr *error, fn func()) {
	if *errPtr != nil {
		fn()
	}
}

func OnErrorOrPanic(errPtr *error, fn func()) {
	recoverObj := recover()
	if recoverObj != nil || *errPtr != nil {
		fn()
	}
	if recoverObj != nil {
		panic(recoverObj)
	}
}

func OnPanic(fn func()) {
	recoverObj := recover()
	if recoverObj != nil {
		fn()
		panic(recoverObj)
	}
}

func OnSuccess(errPtr *error, fn func()) {
	recoverObj := recover()
	if recoverObj == nil && *errPtr == nil {
		fn()
	}
	if recoverObj != nil {
		panic(recoverObj)
	}
}
