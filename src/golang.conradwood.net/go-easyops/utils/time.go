package utils

import (
	"fmt"
	"time"
)

// format a timestring
func TimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// format a timestamp
func TimestampString(timestamp uint32) string {
	if timestamp == 0 {
		return "unset"
	}
	t := time.Unix(int64(timestamp), 0)
	return TimeString(t)
}

// format a timestamp as 'age'
func TimestampAgeString(timestamp uint32) string {
	if timestamp == 0 {
		return "not set"
	}
	secs_age := uint32(time.Now().Unix()) - timestamp
	if secs_age == 0 {
		return "now"
	}
	minutes := uint32(0)
	secs := secs_age
	if secs_age > 60 {
		minutes = uint32(secs_age / 60)
		secs = secs_age - minutes*60
	}
	var ts []string
	if minutes != 0 {
		ts = append(ts, fmt.Sprintf("%dm", minutes))
	}
	if secs != 0 {
		ts = append(ts, fmt.Sprintf("%ds", secs))
	}
	deli := ""
	res := ""
	for _, s := range ts {
		res = res + deli + s
		deli = " "
	}
	res = res + " ago"
	return res
}
