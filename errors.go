package erx

import "fmt"

func ConditionErr(cond bool, format string, a ...interface{}) error {
	if cond {
		return fmt.Errorf(format, a...)
	}
	return nil
}

func AssertErr(cond bool, format string, a ...interface{}) error {
	if !cond {
		return fmt.Errorf(format, a...)
	}
	return nil
}

func ResultErr(_ interface{}, err error) error {
	return err
}
