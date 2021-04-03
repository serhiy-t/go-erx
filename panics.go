package erx

import "fmt"

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type PanicError struct {
	PanicObj interface{}
}

func (p *PanicError) Error() string {
	return fmt.Sprintf("panic: %v", p.PanicObj)
}

func ErrFromPanic(errPtr *error) {
	r := recover()
	if r != nil {
		rErr, isErr := r.(error)
		if isErr {
			*errPtr = rErr
		} else {
			*errPtr = &PanicError{PanicObj: r}
		}
	}
}
