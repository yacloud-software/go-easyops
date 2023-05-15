package utils

import (
	"fmt"
	"time"
)

var (
	// month=1, day=2
	time_formats = []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02/01/2006",
		"20060102T150405Z",
		"20060102T150405",
		"20060102Z",
		"20060102",
	}
)

// Parse a time in various (uk) formats and return a unix timestamp
func ParseTime(ts string) (uint32, error) {
	for _, tf := range time_formats {
		t, err := time.Parse(tf, ts)
		if err != nil {
			//	fmt.Printf("%s is not formatted to \"%s\": %s\n", ts, tf, err)
			continue
		}
		return uint32(t.Unix()), nil
	}
	return 0, fmt.Errorf("unknown time format \"%s\"", ts)
}
