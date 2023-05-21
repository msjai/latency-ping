package lib

import "time"

// GetDurationFromPtr -
func GetDurationFromPtr(val *time.Duration) time.Duration {
	if val == nil {
		return time.Duration(0)
	}

	return *val
}
