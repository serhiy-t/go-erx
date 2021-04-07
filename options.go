package erx

import (
	"runtime/debug"
)

type Option interface{
	reportIgnored(error)
	reportSuppressed(error)
}

type ErrLogger struct {
	errPtr *error

	ignored    []error
	suppressed []error
}

func ErrorLogger(errPtr *error) *ErrLogger {
	return &ErrLogger{
		errPtr: errPtr,
	}
}

func (e *ErrLogger) LogSuppressedErrors() {
	e.log(e.selectSuppressed())
}

func (e *ErrLogger) LogSilentErrors() {
	e.log(e.selectSuppressed(), e.selectIgnored())
}

func (e *ErrLogger) LogIgnoredErrors() {
	e.log(e.selectIgnored())
}

func (e *ErrLogger) LogAllErrors() {
	e.log(e.selectReturned(), e.selectSuppressed(), e.selectIgnored())
}

func (e *ErrLogger) LogReturnedError() {
	e.log(e.selectReturned())
}

type ErrorForLog struct {
	err error
	tag string
}

type ErrorSelector func(e *ErrLogger) []ErrorForLog

func (e *ErrLogger) reportIgnored(err error) {
	e.ignored = append(e.ignored, err)
}

func (e *ErrLogger) reportSuppressed(err error) {
	e.suppressed = append(e.suppressed, err)
}

func (e *ErrLogger) log(selectors ...ErrorSelector) {
	var stack string
	for _, selector := range selectors {
		errors := selector(e)
		for _, err := range errors {
			if len(stack) == 0 {
				stack = string(debug.Stack())
			}
			globalLogFn(&LogMessage{
				Format: "error: %v",
				A:      []interface{}{err.err},
				Stack:  func() string { return stack },
				Tags:   []string{"erx", err.tag},
			})
		}
	}
}

func (e *ErrLogger) selectReturned() ErrorSelector {
	return func(e *ErrLogger) []ErrorForLog {
		if *e.errPtr == nil {
			return nil
		}
		return []ErrorForLog{
			{
				err: *e.errPtr,
				tag: "returned",
			},
		}
	}
}

func (e *ErrLogger) selectIgnored() ErrorSelector {
	return func(e *ErrLogger) []ErrorForLog {
		var result []ErrorForLog
		for _, err := range e.ignored {
			result = append(result, ErrorForLog{
				err: err,
				tag: "ignored",
			})
		}
		return result
	}
}

func (e *ErrLogger) selectSuppressed() ErrorSelector {
	return func(e *ErrLogger) []ErrorForLog {
		var result []ErrorForLog
		for _, err := range e.suppressed {
			result = append(result, ErrorForLog{
				err: err,
				tag: "suppressed",
			})
		}
		return result
	}
}
