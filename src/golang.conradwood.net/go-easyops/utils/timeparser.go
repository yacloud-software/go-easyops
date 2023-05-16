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
func ParseTime(ts string) (uint32, error) {
	d, err := strconv.ParseUint(ts, 10, 64)
	if err == nil {
		// it's just a number, parse as timestamp
		return uint32(d), nil
	}
	for _, tf := range time_formats {
		t, err := time.Parse(tf, ts)
		if err != nil {
			//	fmt.Printf("%s is not formatted to \"%s\": %s\n", ts, tf, err)
			continue
		}
		return uint32(t.UTC().Unix()), nil
	}
	return 0, fmt.Errorf("unknown time format \"%s\"", ts)
}

// Parse a time in various (uk) formats and return a unix timestamp (in UTC)
func ParseTimeWithLocation(ts, tz string) (uint32, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return 0, err
	}
	for _, tf := range time_formats {
		t, err := time.ParseInLocation(tf, ts, loc)
		if err != nil {
			//	fmt.Printf("%s is not formatted to \"%s\": %s\n", ts, tf, err)
			continue
		}
		return uint32(t.UTC().Unix()), nil
	}
	return 0, fmt.Errorf("unknown time format \"%s\"", ts)
}
