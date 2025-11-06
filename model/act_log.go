package model

import (
	"fmt"
	"strings"
	"time"
)

type ActLog struct {
	MessageID int64
	UserID    int64
	Name      string
	StartTime time.Time
	EndTime   time.Time
}

func (a *ActLog) Text() string {
	endTime := "еще идет"
	if !a.EndTime.IsZero() {
		d := a.EndTime.Sub(a.StartTime)
		ds := d.Truncate(time.Second).String()
		if d >= time.Minute {
			ds = strings.TrimSuffix(d.Truncate(time.Minute).String(), "0s")
		}
		endTime = fmt.Sprintf("%s (%s)", a.EndTime.Format("02/01/2006 15:04"), ds)
	}

	return fmt.Sprintf("Занятие: %s\nНачало: %s\nКонец: %s",
		a.Name, a.StartTime.Format("02/01/2006 15:04"), endTime)
}
