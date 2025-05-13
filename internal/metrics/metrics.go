package metrics

import "time"

// Returns the elapsed time in milliseconds.
func Elapsed(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}
