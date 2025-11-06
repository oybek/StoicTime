package service

import (
	"time"
)

var DefaultTZ = time.FixedZone("BSK", 6*60*60)

type Clock struct{}

func (c *Clock) Now() time.Time {
	return time.Now().UTC()
}

func (c *Clock) TodayMidnight(loc *time.Location) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).UTC()
}

func (c *Clock) TomorrowMidnight(loc *time.Location) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc).UTC()
}
