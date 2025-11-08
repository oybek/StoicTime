package service

import (
	"time"
)

var DefaultTZ = time.FixedZone("BSK", 6*60*60)

type Clock struct{}

func (c *Clock) Now() time.Time {
	return time.Now().UTC()
}

func (c *Clock) TodayStartEnd(loc *time.Location) (time.Time, time.Time) {
	now := time.Now().In(loc)
	startTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, loc,
	)
	endTime := startTime.Add(24 * time.Hour)
	return startTime.UTC(), endTime.UTC()
}
