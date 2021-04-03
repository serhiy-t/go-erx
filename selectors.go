package erx

import "runtime/debug"

func (e *ErrorReference) LogSuppressedErrors() {
	e.log(e.selectSuppressed())
}

func (e *ErrorReference) LogSilentErrors() {
	e.log(e.selectSuppressed(), e.selectIgnored())
}

func (e *ErrorReference) LogIgnoredErrors() {
	e.log(e.selectIgnored())
}

func (e *ErrorReference) LogAllErrors() {
	e.log(e.selectReturned(), e.selectSuppressed(), e.selectIgnored())
}

func (e *ErrorReference) LogReturnedError() {
	e.log(e.selectReturned())
}

type ErrorForLog struct {
	err error
	tag string
}

type ErrorSelector func(e *ErrorReference) []ErrorForLog

func (e *ErrorReference) log(selectors ...ErrorSelector) {
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

func (e *ErrorReference) selectReturned() ErrorSelector {
	return func(e *ErrorReference) []ErrorForLog {
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

func (e *ErrorReference) selectIgnored() ErrorSelector {
	return func(e *ErrorReference) []ErrorForLog {
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

func (e *ErrorReference) selectSuppressed() ErrorSelector {
	return func(e *ErrorReference) []ErrorForLog {
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
