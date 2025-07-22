package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type RepeatingNotificationJob struct {
	Interval  time.Duration
	Cron      *cron.Cron
	EntryID   *cron.EntryID
	Scheduler Scheduler
}

func (j RepeatingNotificationJob) Run() {
	entryID, _ := j.Cron.AddFunc(
		fmt.Sprintf("@every %s", j.Interval), j.Scheduler.Notifier,
	)

	*j.EntryID = entryID
	fmt.Printf("▶️ Начал уведомления (%s) Задача %d\n", j.Interval, j.EntryID)
}

type StopNotificationJob struct {
	Cron       *cron.Cron
	EntryID *cron.EntryID
}

func (j StopNotificationJob) Run() {
	j.Cron.Remove(*j.EntryID)
	fmt.Printf("⛔ Уведомления остановлены (%s)\n", time.Now())
}
