package erx

import "fmt"

type ErrorReference struct {
	errPtr *error

	ignored    []error
	suppressed []error
}

type RetFn func(err error, errs ...error) bool

func (e *ErrorReference) ret(err error, errs ...error) bool {
	found := false
	if GetActualError(err) != nil {
		*e.errPtr = GetActualError(err)
		found = true
	}
	for _, err = range errs {
		if GetActualError(err) != nil {
			if !found {
				*e.errPtr = GetActualError(err)
				found = true
			} else {
				e.suppressed = append(e.suppressed, GetActualError(err))
			}
		}
	}
	return found
}

func ErrorRef(err *error) (*ErrorReference, RetFn) {
	result := &ErrorReference{errPtr: err}
	return result, result.ret
}

func (e *ErrorReference) GetActualError() error {
	return *e.errPtr
}

func (e *ErrorReference) Error() string {
	err := e.GetActualError()
	if err == nil {
		return "<nil>"
	}
	return err.Error()
}

func (e *ErrorReference) Wrap(fn func(error) error) {
	Wrap(e.errPtr, fn)
}

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

func (e *ErrorReference) WrapFmtErrorf(format string, a ...interface{}) {
	WrapFmtErrorf(e.errPtr, format, a...)
}

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

func (e *ErrorReference) WrapFmtErrorw(format string, a ...interface{}) {
	WrapFmtErrorw(e.errPtr, format, a...)
}

func WrapFmtErrorw(errPtr *error, format string, a ...interface{}) {
	if *errPtr == nil {
		return
	}
	WrapFmtErrorf(errPtr, format+": %w", append(a, OriginalErr)...)
}

func (e *ErrorReference) CheckErr(fn func() error) {
	fnErr := fn()
	if fnErr != nil {
		if *e.errPtr == nil {
			*e.errPtr = fnErr
		} else {
			e.suppressed = append(e.suppressed, fnErr)
		}
	}
}

func CheckErr(errPtr *error, fn func() error) {
	fnErr := fn()
	if fnErr != nil && *errPtr == nil {
		*errPtr = fnErr
	}
}

func (e *ErrorReference) IgnoreErr(fn func() error) {
	ignoredErr := fn()
	if ignoredErr != nil {
		e.ignored = append(e.ignored, ignoredErr)
	}
}

func IgnoreErr(fn func() error) {
	_ = fn()
}

func (e *ErrorReference) OnError(fn func()) {
	OnError(e.errPtr, fn)
}

func OnError(errPtr *error, fn func()) {
	if *errPtr != nil {
		fn()
	}
}

func (e *ErrorReference) OnErrorOrPanic(fn func()) {
	onErrorOrPanic(recover(), e.errPtr, fn)
}

func OnErrorOrPanic(errPtr *error, fn func()) {
	onErrorOrPanic(recover(), errPtr, fn)
}

func onErrorOrPanic(recoverObj interface{}, errPtr *error, fn func()) {
	if recoverObj != nil || *errPtr != nil {
		fn()
	}
	if recoverObj != nil {
		panic(recoverObj)
	}
}

func (e *ErrorReference) OnPanic(fn func()) {
	onPanic(recover(), fn)
}

func OnPanic(fn func()) {
	onPanic(recover(), fn)
}

func onPanic(recoverObj interface{}, fn func()) {
	if recoverObj != nil {
		fn()
		panic(recoverObj)
	}
}

func (e *ErrorReference) OnSuccess(fn func()) {
	onSuccess(recover(), e.errPtr, fn)
}

func OnSuccess(errPtr *error, fn func()) {
	onSuccess(recover(), errPtr, fn)
}

func onSuccess(recoverObj interface{}, errPtr *error, fn func()) {
	if recoverObj == nil && *errPtr == nil {
		fn()
	}
	if recoverObj != nil {
		panic(recoverObj)
	}
}
