package timeUtils

import "time"

// CleanTime truncates time(hours, seconds, ...) from t
func CleanTime(t time.Time) time.Time {
	return t.Truncate(time.Hour).Add(time.Duration(t.Hour()) * time.Hour)
}
