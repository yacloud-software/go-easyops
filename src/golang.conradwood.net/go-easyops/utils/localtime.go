package utils

import (
	"context"
	"time"
)

func LocalTime(ctx context.Context) (time.Time, error) {
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		return time.Now(), err
	}
	now := time.Now().In(loc)
	return now, nil
}
