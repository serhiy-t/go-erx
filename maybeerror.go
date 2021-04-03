package erx

type MaybeError interface {
	GetActualError() error
}

func GetActualError(err error) error {
	maybeErr, ok := err.(MaybeError)
	if ok {
		return maybeErr.GetActualError()
	}
	return err
}
