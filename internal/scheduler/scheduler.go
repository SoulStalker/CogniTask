package scheduler

import (
	"log"
	"time"

	"github.com/SoulStalker/cognitask/internal/handlers"
	"github.com/SoulStalker/cognitask/internal/messages"
	"github.com/SoulStalker/cognitask/internal/usecase"
	
	tele "gopkg.in/telebot.v3"
)

type Notifier struct {
	settingsSvc usecase.SettingsService
	tasksSvc    usecase.TaskService
	ch          chan (time.Duration)
	bot         *tele.Bot
}

func NewNotifier(
	settingsSvc usecase.SettingsService,
	tasksSvc usecase.TaskService,
	ch chan (time.Duration),
	bot *tele.Bot,
) *Notifier {
	return &Notifier{
		settingsSvc: settingsSvc,
		tasksSvc:    tasksSvc,
		ch:          ch,
		bot:         bot,
	}
}

// TaskNotificationsScheduler -  отправляет список задач по каждый указанный интервал в часах
func (n *Notifier) TaskNotificationsScheduler(chaiID int64) {
	recipient := tele.ChatID(chaiID)

	newInterval, err := n.settingsSvc.GetNotificationInterval()
	if err != nil {
		log.Println(err.Error())
	}

	ticker := time.NewTicker(time.Duration(newInterval) * time.Hour)

	for {
		select {
		case <-ticker.C:
			currentTasks, err := n.tasksSvc.GetPending()
			if err != nil {
				n.bot.Send(recipient, "Не смог отправить задачи:\n"+err.Error())
			}

			rows := handlers.FormatTaskList(currentTasks)
			res, err := n.bot.Send(recipient, messages.BotMessages.YourTasks, &tele.ReplyMarkup{InlineKeyboard: rows})

			if err != nil {
				log.Printf("Не смог отправить задачи:\n%v, %v", err, res)
			}

		case interval := <-n.ch:
			ticker.Stop()
			ticker = time.NewTicker(interval * time.Hour)
		}
	}
}
