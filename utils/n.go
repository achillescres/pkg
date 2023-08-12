package utils

func TryNTimes[R any](n int, f func() (R, error)) (R, error) {
	var r R
	var err error
	for i := 0; i < n; i++ {
		r, err = f()
		if err == nil {
			break
		}
	}
	return r, err
}
