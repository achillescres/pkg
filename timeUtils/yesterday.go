package timeUtils

import "time"

func Yesterday(t time.Time) time.Time {
	return CleanTime(t).AddDate(0, 0, -1)
}
