package util

import "time"

// NowUTC returns the current time in UTC
func NowUTC() time.Time {
	return time.Now().UTC().Round(time.Millisecond)
}
