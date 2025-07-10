package scheduler

import (
	"fmt"
	"time"
)

// TaskNotificationsScheduler -  отправляет список задач по каждый указанный интервал в часах

func TaskNotificationsScheduler(newInterval chan(time.Duration)) {
	ticker := time.NewTicker(<-newInterval)

	for {
		select {
		case <- ticker.C:
			fmt.Println("🔔 Напоминание! Время:", time.Now())
		case interval := <-newInterval:
			fmt.Println("⏱ Новый интервал:", interval)
			ticker.Stop()
			ticker = time.NewTicker(interval)
		}
	}
}