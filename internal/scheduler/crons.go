package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/SoulStalker/cognitask/internal/handlers"
	"github.com/SoulStalker/cognitask/internal/messages"
	"github.com/SoulStalker/cognitask/internal/usecase"
	"github.com/robfig/cron/v3"
	tele "gopkg.in/telebot.v3"
)

// надо добавить расписания для автоудаления сообщений (проверять раз в день и удалять все что старше заданных дней)
// расписание для моитиваций
// расписание для запуска уведомлений
// расписание для остановки уведомлений
// либо перенсти тикер в крон

type Scheduler struct {
	cr         *cron.Cron
	settingsUC *usecase.SettingsService
	taskUC     *usecase.TaskService
	mediaUC    *usecase.MediaService
	bot        *tele.Bot
	chatID     int64
}

func NewScheduler(
	cr *cron.Cron,
	settingsUC *usecase.SettingsService,
	taskUC *usecase.TaskService,
	mediaUC *usecase.MediaService,
	bot *tele.Bot,
	chatID int64,
) *Scheduler {
	return &Scheduler{
		cr:         cr,
		settingsUC: settingsUC,
		taskUC:     taskUC,
		mediaUC:    mediaUC,
		bot:        bot,
		chatID:     chatID,
	}
}

// InitDefaultSchedule запускает расписания по настройками из базы
func (s *Scheduler) InitDefaultSchedule() {
	// авто-удаление старых задач
	s.deleteOldTasks()
	// расписаний мотиваций замокано
	s.SendMediaJob(14)
	// можно начинать слать задачи с интервалом
	var notifyEntryID cron.EntryID

	startJob := RepeatingNotificationJob{
		Interval: 3 * time.Hour,
		Cron:     s.cr,
		EntryID:  &notifyEntryID,
		Scheduler: *s,
	}

	stopJob := StopNotificationJob{
		Cron:       s.cr,
		EntryIDPtr: &notifyEntryID,
	}

	s.cr.AddJob("00 9 * * *", startJob)
	s.cr.AddJob("18 16 * * *", stopJob)

	now := time.Now()

	start := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.Local)

	// если бот перезапустился после запуска расписания, запускаем задание
	if now.After(start) && now.Before(end) {
		startJob.Run()
	}

	s.cr.Start()
}

// Задача по удаление старых записей из базы
func (s *Scheduler) deleteOldTasks() {
	N_days, err := s.settingsUC.GetExpirationDays()
	if err != nil {
		log.Panicln(err)
	}
	// Старые записи удаляются хардкодным расписанием в 4:30 утра
	s.cr.AddFunc("30 4 * * *", func() { s.taskUC.RemoveOldTasks(int(N_days)) })
}

// Задача отправляет медиа файлы из базы в заданное время
func (s *Scheduler) SendRandomMedia() {
	recipient := tele.ChatID(s.chatID)

	media, err := s.mediaUC.Random()
	switch media.Type {
	case "photo":
		photo := &tele.Photo{File: tele.File{FileID: media.Link}}
		_, err = s.bot.Send(recipient, photo)
	case "video":
		video := &tele.Video{File: tele.File{FileID: media.Link}}
		_, err = s.bot.Send(recipient, video)
	}
	if err != nil {
		s.bot.Send(recipient, err.Error())
	}
}

func (s *Scheduler) SendMediaJob(hour uint) {
	s.cr.AddFunc(fmt.Sprintf("10 %d * * *", hour), s.SendRandomMedia)
}

func (s *Scheduler) Notifier() {
	recipient := tele.ChatID(s.chatID)

	currentTasks, err := s.taskUC.GetPending()
	if err != nil {
		s.bot.Send(recipient, "Не смог отправить задачи:\n"+err.Error())
	}

	rows := handlers.FormatTaskList(currentTasks)
	res, err := s.bot.Send(recipient, messages.BotMessages.YourTasks, &tele.ReplyMarkup{InlineKeyboard: rows})

	if err != nil {
		log.Printf("Не смог отправить задачи:\n%v, %v", err, res)
	}

}
