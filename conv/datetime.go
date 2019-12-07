package conv

import (
	"time"
)

func StrToDateTime(format, str string, local *time.Location) time.Time {
	location, _ := time.ParseInLocation(format, str, local)
	return location
}

func StrToDateTimePtr(format, str string, local *time.Location) *time.Time {
	if str == "" {
		return nil
	}
	location, _ := time.ParseInLocation(format, str, local)
	return &location
}

func DateTimeToStr(format string, t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(format)
}
