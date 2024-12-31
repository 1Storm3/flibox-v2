package mapper

import (
	"time"
)

func parseTimeStringToTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}
