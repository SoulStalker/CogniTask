package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/SoulStalker/cognitask/internal/usecase"
)

type Notifier struct {
	settings usecase.SettingsService
	ch       chan (time.Duration)
}

func NewNotifier(settings usecase.SettingsService, ch       chan (time.Duration)) *Notifier {
	return &Notifier{
		settings: settings,
		ch: ch,
	}
}

// TaskNotificationsScheduler -  отправляет список задач по каждый указанный интервал в часах
func (n *Notifier) TaskNotificationsScheduler() {
	newInterval, err := n.settings.Interval()
	if err != nil {
		log.Println(err.Error())
	}
	ticker := time.NewTicker(time.Duration(newInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("🔔 Напоминание! Время:", time.Now())
		case interval := <-n.ch:
			fmt.Println("⏱ Новый интервал:", interval)
			ticker.Stop()
			ticker = time.NewTicker(interval)
		}
	}
}
