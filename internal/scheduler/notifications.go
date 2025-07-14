package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type RepeatingNotificationJob struct {
	Interval time.Duration
	Cron     *cron.Cron
	EntryID  *cron.EntryID
}

func (j RepeatingNotificationJob) Run() {
	entryID, _ := j.Cron.AddFunc(
		fmt.Sprintf("@every %s", j.Interval),
		func() {
			fmt.Printf("🔔 Повторное уведомление (%s)\n", time.Now())
		},
	)

	*j.EntryID = entryID
	fmt.Printf("▶️ Начал уведомления (%s)\n", j.Interval)
}

type StopNotificationJob struct {
	Cron       *cron.Cron
	EntryIDPtr *cron.EntryID
}

func (j StopNotificationJob) Run() {
	j.Cron.Remove(*j.EntryIDPtr)
	fmt.Printf("⛔ Уведомления остановлены (%s)\n", time.Now())
}
