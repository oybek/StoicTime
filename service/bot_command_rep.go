package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"oybek.io/sigma/model"
	"oybek.io/sigma/rdb"
)

func (b *Bot) handleCommandRep(tg *gotgbot.Bot, tgctx *ext.Context) error {
	ctx := context.Background()

	todayMidnight := b.clock.TodayMidnight(DefaultTZ)
	tomorrowMidnight := b.clock.TomorrowMidnight(DefaultTZ)
	actLogs, err := b.actLogStorage.FindActLog(ctx, rdb.FindActLogArg{
		UserID:   tgctx.EffectiveChat.Id,
		FromTime: todayMidnight,
		ToTime:   tomorrowMidnight,
	})
	if err != nil {
		return err
	}

	actLogs = trim(actLogs, todayMidnight, tomorrowMidnight)

	s := todayMidnight.In(DefaultTZ).Format("02/01/2006") + "\n\n"
	last := todayMidnight
	for _, al := range actLogs {
		if al.StartTime.Sub(last) > time.Minute {
			s += formatRange(last, al.StartTime, "❓", DefaultTZ) + "\n"
			s += formatRange(al.StartTime, al.EndTime, al.Name, DefaultTZ) + "\n"
			last = al.EndTime
		}
	}
	if tomorrowMidnight.Sub(last) > time.Minute {
		s += formatRange(last, tomorrowMidnight, "❓", DefaultTZ) + "\n"
	}

	_, err = tg.SendMessage(tgctx.EffectiveChat.Id, s+"\n"+summary(actLogs), &gotgbot.SendMessageOpts{
		ParseMode: "markdown",
	})
	return err
}

func summary(logs []model.ActLog) string {
	totalDurations := make(map[string]time.Duration)

	for _, log := range logs {
		duration := log.EndTime.Sub(log.StartTime)
		totalDurations[log.Name] += duration
	}

	type kv struct {
		Name     string
		Duration time.Duration
	}
	var sorted []kv
	for name, dur := range totalDurations {
		sorted = append(sorted, kv{name, dur})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Duration > sorted[j].Duration
	})

	s := ""
	for _, item := range sorted {
		s += fmt.Sprintf("`%s` %v\n", model.FormatDuration(item.Duration), item.Name)
	}
	return s
}

func formatRange(start, end time.Time, name string, loc *time.Location) string {
	return fmt.Sprintf(
		"`%s - %s %6s` %s ",
		start.In(loc).Format("15:04"),
		end.In(loc).Format("15:04"),
		model.FormatDuration(end.Sub(start)),
		name,
	)
}

func trim(als []model.ActLog, A, B time.Time) []model.ActLog {
	if len(als) == 0 {
		return als
	}

	if als[0].StartTime.Before(A) {
		als[0].StartTime = A
	}
	if als[len(als)-1].EndTime.After(B) {
		als[0].EndTime = B
	}

	return als
}
