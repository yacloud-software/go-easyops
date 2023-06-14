package utils

import (
	"fmt"
	"strconv"
	"time"
)

var (
	// month=1, day=2
	time_formats = []string{
		"02/01/2006 15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02/01/2006",
		"20060102T150405Z",
		"20060102T150405",
		"20060102Z",
		"20060102",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05",
	}
)

// Parse a time in various (uk) formats and return a unix timestamp
func ParseTimestamp(ts string) (uint32, error) {
	t, err := ParseTime(ts)
	if err != nil {
		return 0, err
	}
	return uint32(t.Unix()), nil
}

// Parse a time in various (uk) formats and return a unix timestamp
func ParseTimestampWithLocation(ts, tz string) (uint32, error) {
	t, err := ParseTimeWithLocation(ts, tz)
	if err != nil {
		return 0, err
	}
	return uint32(t.Unix()), nil
}

// Parse a time in various (uk) formats and return a time.Time
func ParseTime(ts string) (time.Time, error) {
	d, err := strconv.ParseUint(ts, 10, 64)
	if err == nil {
		// it's just a number, parse as timestamp
		t := time.Unix(int64(d), 0)
		return t, nil
	}
	for _, tf := range time_formats {
		t, err := time.Parse(tf, ts)
		if err != nil {
			//	fmt.Printf("%s is not formatted to \"%s\": %s\n", ts, tf, err)
			continue
		}
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("unknown time format \"%s\"", ts)
}

// Parse a time in various (uk) formats and return a unix timestamp (in UTC)
func ParseTimeWithLocation(ts, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}
	for _, tf := range time_formats {
		t, err := time.ParseInLocation(tf, ts, loc)
		if err != nil {
			//	fmt.Printf("%s is not formatted to \"%s\": %s\n", ts, tf, err)
			continue
		}
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("unknown time format \"%s\"", ts)
}
