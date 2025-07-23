package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/SoulStalker/cognitask/internal/mappers"
	"github.com/SoulStalker/cognitask/internal/messages"
	"github.com/robfig/cron/v3"

	tele "gopkg.in/telebot.v3"
)

type RepeatingNotificationJob struct {
	Interval  time.Duration
	Cron      *cron.Cron
	EntryID   *cron.EntryID
	Scheduler Scheduler
}

func (j *RepeatingNotificationJob) Run() {
	entryID, _ := j.Cron.AddFunc(
		fmt.Sprintf("@every %s", j.Interval), j.Scheduler.Notifier,
	)

	*j.EntryID = entryID
	fmt.Printf("▶️ Начал уведомления (%s) Задача %d\n", j.Interval, j.EntryID)
}

// Notifier отправляет уведомления в телеграмм
func (s *Scheduler) Notifier() {
	if !s.notificationTime() {
		log.Println("Не время беспокоить! ")
		return
	}

	recipient := tele.ChatID(s.chatID)

	currentTasks, err := s.taskUC.GetPending()
	if err != nil {
		log.Printf("Не смог получить задачи: %v", err)
		s.bot.Send(recipient, "Не смог получить задачи: "+err.Error())
		return
	}

	if len(currentTasks) == 0 {
		s.bot.Send(recipient, "Новых задач нет ")
		return
	}

	rows := mappers.FormatTaskList(currentTasks)
	res, err := s.bot.Send(recipient, messages.BotMessages.YourTasks, &tele.ReplyMarkup{InlineKeyboard: rows})

	if err != nil {
		log.Printf("Не смог отправить задачи:\n%v, %v", err, res)
	}
}
