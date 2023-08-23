package timeUtils

import "time"

// Yesterday returns yesterday time with cleaned time
func Yesterday(t time.Time) time.Time {
	return CleanTime(t).AddDate(0, 0, -1)
}
