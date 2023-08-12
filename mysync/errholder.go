package mysync

// ErrorHolder holds first non-nil error forever
type ErrorHolder struct {
	err error
}

func (eH *ErrorHolder) Hold(err error) {
	if eH.err != nil {
		return
	}
	eH.err = err
}

func (eH *ErrorHolder) Error() error {
	return eH.err
}
